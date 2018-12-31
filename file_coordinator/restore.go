package file_coordinator

import (
	"errors"
	"fmt"
	"github.com/luhring/awsmfa/credentials_file"
	"github.com/luhring/awsmfa/environment"
)

func (c *Coordinator) Restore() error {
	if c.Env.DoesHaveCredentialsFile() {
		credentialsFile, err := c.getCredentialsFile()

		if err != nil {
			return err
		}

		if credentialsFile.DoesProfileHavePermanentCredentials(c.SelectedProfileName) {
			fmt.Printf("'%s' profile already contains permanent credentials\n", c.SelectedProfileName)
			return nil
		}

		// credentials file exists but has temporary credentials, so there's no longer a use for this file

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

		areCredentialsTemporary := false == credentialsFile.DoesProfileHavePermanentCredentials(c.SelectedProfileName)
		willCredentialsFileBeUsed := false == environment.WillEnvironmentVariablesPreemptUseOfCredentialsFile()

		if areCredentialsTemporary && willCredentialsFileBeUsed && c.Env.DoesHaveCredentialsFileBackup() {
			_ = c.Restore()
		}
	}
}
