package app

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// State represents local application state and reflects a current status
// and a results of a previous run.
type State struct {
	Users   []string `json:"users"`
	Sudoers []string `json:"sudoers"`
}

// Update runtime state.
//
// **Warning**: This will not save the state to disk.
func (s *State) Update(users, sudoers []string) {
	s.Users = users
	s.Sudoers = sudoers
}

// Save runtime state to disk
func (s *State) Save(file string) error {
	stateFile, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "error opening state file")
	}
	defer stateFile.Close()

	content, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "error marshaling state to json")
	}

	if err = ioutil.WriteFile(file, content, 0644); err != nil {
		return errors.Wrap(err, "error writing state file")
	}

	return nil
}

// LoadState from disk
func LoadState(file string) (*State, error) {
	stateFile, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "error opening state file")
	}
	defer stateFile.Close()

	content, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading state file")
	}

	o := new(State)
	if err = json.Unmarshal(content, o); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling state file")
	}

	return o, nil
}
