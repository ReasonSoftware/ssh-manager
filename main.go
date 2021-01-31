package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ReasonSoftware/ssh-manager/internal/app"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// Version of an application
	Version string = "1.0.0"
	// AppDir contains an home dir for an application files
	AppDir string = "/var/lib/ssh-manager"
	// StateFile contains a filename of a state file
	StateFile string = "state.json"
)

func init() {
	// logger
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	// config
	viper.SetConfigName("ssh-manager")
	viper.SetConfigType("yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/root")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "error reading configuration file"))
	}

	if len(viper.GetStringSlice("groups")) == 0 {
		log.Fatal("configuration does not contain any groups")
	}

	// state
	_, err := os.Stat(AppDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(AppDir, 0777); err != nil {
			log.Fatal(errors.Wrap(err, "error creating application directory"))
		}
	} else if err != nil {
		log.Fatal(errors.Wrap(err, "error validating application directory"))
	}

	_, err = os.Stat(path.Join(AppDir, StateFile))
	if os.IsNotExist(err) {
		if err = ioutil.WriteFile(path.Join(AppDir, StateFile), []byte("{}"), 0666); err != nil {
			log.Fatal(errors.Wrap(err, "error creating state file"))
		}
	} else if err != nil {
		log.Fatal(errors.Wrap(err, "error validating state file"))
	}

	// groups
	if err := app.EnsureGroups(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Infof("ssh-manager v%v started", Version)

	state, err := app.LoadState(path.Join(AppDir, StateFile))
	if err != nil {
		log.Fatal(errors.Wrap(err, "error loading state"))
	}
	log.Info("configured server groups: ", app.SliceToString(viper.GetStringSlice("groups")))

	// get members
	secretsManager := secretsmanager.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))

	log.Info("fetching remote configuration")
	conf, err := app.GetConfig(secretsManager, viper.GetString("secret_name"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "error fetching remote configuration"))
	}

	// get unique members
	users := conf.GetUsers(viper.GetStringSlice("groups"))
	listOfUsers := []string{}
	for username := range users {
		listOfUsers = append(listOfUsers, username)
	}
	log.Info("configured users: ", app.SliceToString(listOfUsers))

	sudoers := conf.GetSudoers(viper.GetStringSlice("groups"))
	listOfSudoers := []string{}
	for username := range sudoers {
		listOfSudoers = append(listOfSudoers, username)
	}
	log.Info("configured sudoers: ", app.SliceToString(listOfSudoers))

	// operate users
	state.UsersLoop(users)
	state.SudoersLoop(sudoers, listOfUsers)
	state.DeleteUsers(listOfUsers, listOfSudoers)

	// save state
	state.Update(listOfUsers, listOfSudoers)
	if err := state.Save(path.Join(AppDir, StateFile)); err != nil {
		log.Fatal(errors.Wrap(err, "error saving the state"))
	}
}
