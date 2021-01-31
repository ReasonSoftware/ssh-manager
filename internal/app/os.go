package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// SudoersGroup represent a sudoers unix group name
	SudoersGroup string = "ssh-manager-sudoers"
	// UsersGroup represent a users unix group name
	UsersGroup string = "ssh-manager-users"
)

// EnsureGroups validates that both custom users and sudoers groups exists on a server.
// If not, they will be created.
//
// Validation is partial, only by looking on a sudoers file which should contain a custom
// sudoers group. This change is made at the end of the function, after groups creation.
func EnsureGroups() error {
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
			log.Info("user groups exists")
			return nil
		}

		newContent = append(newContent, line)
	}

	log.Info("adding user groups")

	if err := execShellCommand(fmt.Sprintf("groupadd -g 32109 %v", SudoersGroup)); err != nil {
		return errors.Wrap(err, "error creating sudoers group")
	}
	if err := execShellCommand(fmt.Sprintf("groupadd -g 32108 %v", UsersGroup)); err != nil {
		return errors.Wrap(err, "error creating users group")
	}

	newContent = append(newContent, instruction)

	output := strings.Join(newContent, "\n")
	err = ioutil.WriteFile(f, []byte(output), 0440)
	if err != nil {
		return errors.Wrap(err, "error writing to sudoers file")
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
