package file_coordinator

import (
	"errors"
	"github.com/luhring/awsmfa/credentials_file"
	"github.com/luhring/awsmfa/environment"
)

type Coordinator struct {
	Env                 *environment.Environment
	SelectedProfileName string
}

func New(env *environment.Environment, selectedProfile string) (*Coordinator, error) {
	if env == nil {
		return nil, errors.New("env parameter cannot be nil")
	}

	if len(selectedProfile) == 0 {
		return nil, errors.New("selectedProfile parameter cannot have zero length")
	}

	return &Coordinator{
		env,
		selectedProfile,
	}, nil
}

func (c *Coordinator) getCredentialsFile() (*credentials_file.CredentialsFile, error) {
	if false == c.Env.DoesHaveCredentialsFile() {
		return nil, errors.New("unable to find credentials file")
	}

	return credentials_file.NewFromDisk(c.Env.PathToCredentialsFile())
}

func (c *Coordinator) getCredentialsFileBackup() (*credentials_file.CredentialsFile, error) {
	if false == c.Env.DoesHaveCredentialsFileBackup() {
		return nil, errors.New("unable to find backup of credentials file")
	}

	return credentials_file.NewFromDisk(c.Env.PathToCredentialsFileBackup())
}
