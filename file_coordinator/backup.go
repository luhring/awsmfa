package file_coordinator

import (
	"fmt"
	"github.com/luhring/awsmfa/credentials_file"
)

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

func (c *Coordinator) BackUpPermanentCredentialsIfPresent() error {
	if c.Env.DoesHaveCredentialsFile() {
		credentialsFile, err := c.getCredentialsFile()
		if err != nil {
			return err
		}

		if credentialsFile.DoesProfileHavePermanentCredentials(c.SelectedProfileName) {
			return c.BackUp()
		}
	}

	return nil
}
