package file_coordinator

import (
	"errors"
	"fmt"
	"github.com/luhring/awsmfa/credentials_file"
	"github.com/luhring/awsmfa/environment"
)

type Coordinator struct {
	Env         *environment.Environment
	ProfileName string
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

func (c *Coordinator) BackUp() error {
	if existingBackup, err := c.getCredentialsFileBackup(); err == nil {
		err = existingBackup.Delete()
		if err != nil {
			return err
		}
	}

	credentialsFile, err := c.getCredentialsFile()
	if err != nil {
		return err
	}

	newBackup, err := credentials_file.NewFromConfiguration(credentialsFile.Configuration, c.Env.PathToCredentialsFileBackup())
	if err != nil {
		return err
	}

	err = newBackup.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Backed up credentials file to %s\n", newBackup.Filename)

	return nil
}

func (c *Coordinator) Restore() error {
	if c.Env.DoesHaveCredentialsFile() {
		credentialsFile, err := c.getCredentialsFile()

		if err != nil {
			return err
		}

		if credentialsFile.DoesProfileHavePermanentCredentials(c.ProfileName) {
			fmt.Printf("'%s' profile already contains permanent credentials\n", c.ProfileName)

			return nil
		}

		// credentials file exists but has temporary credentials -- there's no longer a use for this file

		err = credentialsFile.Delete()
		if err != nil {
			return err
		}
	}

	// credentials file doesn't exist

	if c.Env.DoesHaveCredentialsFileBackup() {
		backup, err := c.getCredentialsFileBackup()
		if err != nil {
			return err
		}

		newCredentialsFile, err := credentials_file.NewFromConfiguration(backup.Configuration, c.Env.PathToCredentialsFile())
		if err != nil {
			return err
		}

		err = newCredentialsFile.Save()
		if err != nil {
			return err
		}

		fmt.Println("Restored original credentials from backup")
		fmt.Println("You can no longer perform actions that require MFA")

		return backup.Delete()
	}

	return errors.New("unable to find original credentials")
}

func (c *Coordinator) RestorePermanentCredentialsIfAppropriate() {
	// We can't get a session token using temporary credentials.
	// If we're going to create our AWS client session using the credentials file, it needs to have permanent credentials.
	// As a convenience to the user, we'll detect this scenario and attempt to restore a backup if one exists.
	// Since this behavior isn't part of the critical path, we'll return silently rather than elevate errors.

	if c.Env.DoesHaveCredentialsFile() {
		credentialsFile, err := c.getCredentialsFile()
		if err != nil {
			return
		}

		areCredentialsTemporary := false == credentialsFile.DoesProfileHavePermanentCredentials(c.ProfileName)
		willCredentialsFileBeUsed := false == environment.WillEnvironmentVariablesPreemptUseOfCredentialsFile()

		if areCredentialsTemporary && willCredentialsFileBeUsed && c.Env.DoesHaveCredentialsFileBackup() {
			_ = c.Restore()
		}
	}
}

func (c *Coordinator) BackUpPermanentCredentialsIfPresent() error {
	if c.Env.DoesHaveCredentialsFile() {
		credentialsFile, err := c.getCredentialsFile()
		if err != nil {
			return err
		}

		if credentialsFile.DoesProfileHavePermanentCredentials(c.ProfileName) {
			return c.BackUp()
		}
	}

	return nil
}

func (c *Coordinator) getCredentialsFile() (*credentials_file.CredentialsFile, error) {
	if false == c.Env.DoesHaveCredentialsFile() {
		return nil, errors.New("unable to find credentials file")
	}

	return credentials_file.LoadFromDisk(c.Env.PathToCredentialsFile())
}

func (c *Coordinator) getCredentialsFileBackup() (*credentials_file.CredentialsFile, error) {
	if false == c.Env.DoesHaveCredentialsFileBackup() {
		return nil, errors.New("unable to find backup of credentials file")
	}

	return credentials_file.LoadFromDisk(c.Env.PathToCredentialsFileBackup())
}
