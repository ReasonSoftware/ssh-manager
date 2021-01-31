package app

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// UsersLoop is a main loop for standard users creation and sudoers demotion
func (s *State) UsersLoop(users map[string]string) {
USERS:
	for user, key := range users {
		for _, stateUser := range s.Users {
			if user == stateUser {
				continue USERS
			}
		}

		for _, stateSudoer := range s.Sudoers {
			if user == stateSudoer {
				log.Infof("demoting a user: %v", user)
				if err := DemoteUser(user); err != nil {
					log.Error(errors.Wrapf(err, "error demoting a user %v", user))
				}
				continue USERS
			}
		}

		if err := CreateUsers(user, key, false); err != nil {
			log.Error(errors.Wrapf(err, "error creating a user '%v'", user))
		}
	}
}

// SudoersLoop is a main loop for sudo users creation and standard users promotion
func (s *State) SudoersLoop(sudoers map[string]string, listOfUsers []string) {
SUDOERS:
	for sudoer, key := range sudoers {
		for _, user := range listOfUsers {
			if sudoer == user {
				log.Errorf("user %v promotion denied because of a privilege conflict", sudoer)
				continue SUDOERS
			}
		}

		for _, stateSudoer := range s.Sudoers {
			if sudoer == stateSudoer {
				continue SUDOERS
			}
		}

		for _, stateUser := range s.Users {
			if sudoer == stateUser {
				log.Infof("promoting a user: %v", sudoer)
				if err := PromoteUser(sudoer); err != nil {
					log.Error(errors.Wrapf(err, "error promoting a user %v", sudoer))
				}
				continue SUDOERS
			}
		}

		if err := CreateUsers(sudoer, key, true); err != nil {
			log.Error(errors.Wrapf(err, "error creating a user '%v'", sudoer))
		}
	}
}
