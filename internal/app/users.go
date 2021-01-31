package app

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateUsers will create a local users and assign them to the relevant
// groups.
func CreateUsers(user, key string, sudoer bool) error {
	if sudoer {
		log.Infof("adding a sudoer: %v", user)
	} else {
		log.Infof("adding a user: %v", user)
	}

	// create user
	command := fmt.Sprintf("useradd --create-home --home-dir /home/%v --shell /bin/bash --password %v %v", user, genPassword(), user)
	if err := execShellCommand(command); err != nil {
		return err
	}

	// add public ssh key
	if err := os.MkdirAll(fmt.Sprintf("/home/%v/.ssh", user), 0700); err != nil {
		return errors.Wrap(err, "error creating application directory")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("/home/%v/.ssh/authorized_keys", user), []byte(key), 0600); err != nil {
		return errors.Wrap(err, "error creating authorized_keys file")
	}

	// promote user to sudoer
	if sudoer {
		if err := PromoteUser(user); err != nil {
			return errors.Wrap(err, "error promoting a user")
		}
	} else {
		if err := DemoteUser(user); err != nil {
			return errors.Wrap(err, "error demoting a user")
		}
	}

	return nil
}

func genPassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// DeleteUsers that exists in a runtime state but not in a provided slices,
// which provided from a remote configuration.
func (s *State) DeleteUsers(users, sudoers []string) {
	previousUsers := append(s.Users, s.Sudoers...)
	currentUsers := append(users, sudoers...)
	candidates := make([]string, 0)

	for _, user := range previousUsers {
		removed := true
		for _, u := range currentUsers {
			if user == u {
				removed = false
			}
		}

		if removed {
			candidates = append(candidates, user)
		}
	}

	for _, user := range candidates {
		log.Warnf("removing user: %v", user)

		command := fmt.Sprintf("userdel -r %v", user)
		if err := execShellCommand(command); err != nil {
			log.Errorf(errors.Wrap(err, "error deleting a user").Error())
		}
	}
}

// PromoteUser make standard user a sudo user
func PromoteUser(user string) error {
	command := fmt.Sprintf("usermod -G %v %v", SudoersGroup, user)
	if err := execShellCommand(command); err != nil {
		return err
	}

	return nil
}

// DemoteUser make sudo user a standard user
func DemoteUser(user string) error {
	command := fmt.Sprintf("usermod -G %v %v", UsersGroup, user)
	if err := execShellCommand(command); err != nil {
		return err
	}

	return nil
}
