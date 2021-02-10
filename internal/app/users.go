package app

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"

	"path"
	"strconv"

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

	// add ssh key
	if err := updateAuthorizedKeys(user, key); err != nil {
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

func updateAuthorizedKeys(username, key string) error {
	// get uid/gui
	u, err := user.Lookup(username)
	if err != nil {
		return errors.Wrap(err, "error identifying a user")
	}

	uid, err := strconv.ParseInt(u.Uid, 10, 32)
	if err != nil {
		return errors.Wrap(err, "error identifying a user")
	}

	g, err := user.LookupGroup(username)
	if err != nil {
		return errors.Wrap(err, "error identifying user's group")
	}

	gid, err := strconv.ParseInt(g.Gid, 10, 32)
	if err != nil {
		return errors.Wrap(err, "error identifying user's group")
	}

	dir := path.Join("/home", username, ".ssh")
	file := path.Join(dir, "authorized_keys")

	// create user's .ssh directory
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return errors.Wrap(err, "error creating .ssh directory")
		}

		if err := os.Chown(dir, int(uid), int(gid)); err != nil {
			return errors.Wrap(err, "error updating .ssh ownership")
		}
	} else if err != nil {
		return errors.Wrap(err, "error validating .ssh directory")
	}

	// create user's authorized_keys file
	if err := ioutil.WriteFile(file, []byte(key), 0600); err != nil {
		return errors.Wrap(err, "error writing a file")
	}

	if err := os.Chown(file, int(uid), int(gid)); err != nil {
		return errors.Wrap(err, "error updating authorized_keys ownership")
	}

	return nil
}

func validatePublicKey(username, key string) (bool, error) {
	f := path.Join("/home", username, ".ssh/authorized_keys")

	file, err := os.Open(f)
	if err != nil {
		return false, errors.Wrap(err, "error opening authorized_keys file")
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return false, errors.Wrap(err, "error reading authorized_keys file")
	}

	if string(content) != key {
		if err := updateAuthorizedKeys(username, key); err != nil {
			return false, errors.Wrap(err, "error updating authorized_keys file")
		}

		return true, nil
	}

	return false, nil
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
