package main

import (
	"io/ioutil"
	"os"
	"path"

	"ssh-manager/internal/app"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
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
	viper.SetConfigName("ssh-manager.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/root")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "error reading configuration file"))
	}

	if len(viper.GetStringSlice("groups")) == 0 {
		log.Fatal("configuration does not contain any groups")
	}

	if viper.GetString("secret_name") == "" {
		log.Fatal("configuration does not contain an aws secret name")
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
}

func main() {
	log.Infof("ssh-manager v%v started", app.Version)

	// validate groups
	log.Info("validating users group")
	if err := app.ValidateUsersGroup(); err != nil {
		log.Fatal(err)
	}

	log.Info("validating sudoers group")
	if err := app.ValidateSudoersGroup(); err != nil {
		log.Fatal(err)
	}

	log.Info("validating sudoers group permission")
	if err := app.ValidateSudoersPermissions(); err != nil {
		log.Fatal(err)
	}

	state, err := app.LoadState(path.Join(AppDir, StateFile))
	if err != nil {
		log.Fatal(errors.Wrap(err, "error loading state"))
	}
	log.Info("configured server groups: ", app.SliceToString(viper.GetStringSlice("groups")))

	// get members
	region := viper.GetString("region")
	if region == "" {
		region = "us-east-1"
	}

	secretsManager := secretsmanager.New(session.Must(session.NewSession(&aws.Config{
		Region: &region,
	})))

	log.Info("fetching remote configuration")
	conf, err := app.GetConfig(secretsManager, viper.GetString("secret_name"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "error fetching remote configuration"))
	}

	// warn about staled groups
	for _, group := range viper.GetStringSlice("groups") {
		if _, val := conf.ServerGroups[group]; !val {
			log.Warnf("group %s does not exists on remote configuration", group)
		}
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

	// configure users
	state.UsersLoop(users)
	state.SudoersLoop(sudoers, listOfUsers)
	state.DeleteUsers(listOfUsers, listOfSudoers)

	// save state
	state.Update(listOfUsers, listOfSudoers)
	if err := state.Save(path.Join(AppDir, StateFile)); err != nil {
		log.Fatal(errors.Wrap(err, "error saving the state"))
	}

	log.Info("ssh-manager finished")
}
