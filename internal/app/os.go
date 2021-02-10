package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/pkg/errors"
)

const (
	// SudoersGroup represent a sudoers unix group name
	SudoersGroup string = "ssh-manager-sudoers"
	// UsersGroup represent a users unix group name
	UsersGroup string = "ssh-manager-users"
)

// ValidateSudoersPermissions ensures that sudoers file contains a custom sudoers group.
func ValidateSudoersPermissions() error {
	f := "/etc/sudoers"
	instruction := fmt.Sprintf("%%%v ALL=(ALL) NOPASSWD: ALL", SudoersGroup)

	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "sudoers file does not exists")
	}

	origContent, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Wrap(err, "error reading sudoers file")
	}

	lines := strings.Split(string(origContent), "\n")

	newContent := make([]string, 0)
	for _, line := range lines {
		if line == instruction {
			return nil
		}

		newContent = append(newContent, line)
	}

	newContent = append(newContent, instruction)

	output := strings.Join(newContent, "\n")
	err = ioutil.WriteFile(f, []byte(output), 0440)
	if err != nil {
		return errors.Wrap(err, "error writing to sudoers file")
	}

	return nil
}

// ValidateUsersGroup ensures that custom users group exists
func ValidateUsersGroup() error {
	if err := createGroup(UsersGroup, 32108); err != nil {
		return errors.Wrapf(err, "error validating group %v", UsersGroup)
	}

	return nil
}

// ValidateSudoersGroup ensures that custom sudoers group exists
func ValidateSudoersGroup() error {
	if err := createGroup(SudoersGroup, 32109); err != nil {
		return errors.Wrapf(err, "error validating group %v", SudoersGroup)
	}

	return nil
}

func createGroup(name string, id int64) error {
	_, err := user.LookupGroup(name)
	if err != nil {
		_, unknown := err.(user.UnknownGroupError)

		if unknown {
			if err := execShellCommand(fmt.Sprintf("groupadd -g %v %v", id, name)); err != nil {
				return errors.Wrapf(err, "error creating %v group", name)
			}
		} else {
			return errors.Wrapf(err, "error look up of a group %v", name)
		}
	}

	return nil
}

func execShellCommand(command string) error {
	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:len(strings.Split(command, " "))]...)

	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, strings.ReplaceAll(out.String(), "\n", ";"))
	}

	return nil
}
