package app

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

// Version contains current application version
const Version string = "1.1.0"

// Config represents a remote configuration
type Config struct {
	Users        map[string]string   `json:"users"`
	ServerGroups map[string]*Members `json:"server_groups"`
}

// Members ia a single server group in a configuration
type Members struct {
	Sudoers []string `json:"sudoers"`
	Users   []string `json:"users"`
}

// GetConfig fetches an AWS Secret and returns an application configuration
func GetConfig(service *secretsmanager.SecretsManager, name string) (*Config, error) {
	result, err := service.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	if result.SecretString == nil {
		return nil, errors.New("empty or a binary secret")
	}

	output := &Config{}
	if err := json.Unmarshal([]byte(*result.SecretString), output); err != nil {
		return nil, errors.Wrap(err, "parsing error")
	}

	return output, nil
}
