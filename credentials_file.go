package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

func prepareCredentialsFileForUse() {
	pathToCredentialsBackupFile := getPathToAwsCredentialsBackupFile()

	if doesCredentialsFileExist() {
		if doesCredentialsFileDefaultProfileContainPermanentCredentials() { // as opposed to temporary credentials
			if doesFileExist(pathToCredentialsBackupFile) {
				removeCredentialsFileBackup()
			}

			return
		}

		// Credentials file default profile has only temporary credentials.
		// We won't be able to use those credentials to get new temporary credentials.
		// Let's see if there's a backup file we can restore to the main credentials file.

		if doesFileExist(pathToCredentialsBackupFile) {
			restoreCredentialsFileFromBackup()
		}

		return
	}

	// We don't have usable credentials in the AWS credentials file.
	// Let's see if there's a backup file we can restore to the main credentials file.

	if doesFileExist(pathToCredentialsBackupFile) {
		restoreCredentialsFileFromBackup()
	}
}

func generateCredentialsFileContent(credentials *sts.Credentials) string {
	return fmt.Sprintf(`[default]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s
`, *credentials.AccessKeyId, *credentials.SecretAccessKey, *credentials.SessionToken)
}

func backupCredentialsFileAndSaveNewCredentialsToDisk(newCredentialsFileContent string) {
	pathToCredentialsFile := getPathToAwsCredentialsFile()

	if doesCredentialsFileExist() {
		backUpCredentialsFile()
	}

	err := ioutil.WriteFile(pathToCredentialsFile, []byte(newCredentialsFileContent), 0600)

	if err != nil {
		exitWithErrorMessage("Unable to save new session credentials to %s: %s\n", pathToCredentialsFile, err.Error())
	}

	fmt.Printf("Authentication successful! Saved new session credentials to %s\n", pathToCredentialsFile)
}

func backUpCredentialsFile() {
	const defaultFailureMessageFormat = "Unable to back up credentials file (%s): %s\n"
	credentialsFileBytes, err := ioutil.ReadFile(getPathToAwsCredentialsFile())

	if err != nil {
		exitWithErrorMessage(defaultFailureMessageFormat, getPathToAwsCredentialsFile(), err.Error())
	}

	err = ioutil.WriteFile(getPathToAwsCredentialsBackupFile(), credentialsFileBytes, 0600)

	if err != nil {
		exitWithErrorMessage(defaultFailureMessageFormat, getPathToAwsCredentialsFile(), err.Error())
	}

	fmt.Printf("Created backup of original credentials at %s\n", getPathToAwsCredentialsBackupFile())
}

func restoreCredentialsFileFromBackup() {
	const defaultFailureMessageFormat = "Unable to restore AWS credentials file from backup: %s\n"

	if doesFileExist(getPathToAwsCredentialsFile()) {
		err := os.Remove(getPathToAwsCredentialsFile())

		if err != nil {
			exitWithErrorMessage(defaultFailureMessageFormat, err.Error())
		}
	}

	err := os.Rename(getPathToAwsCredentialsBackupFile(), getPathToAwsCredentialsFile())

	if err != nil {
		exitWithErrorMessage(defaultFailureMessageFormat, err.Error())
	}

	fmt.Printf("Restored credentials file from backup\n")
}

func removeCredentialsFileBackup() {
	pathToCredentialsBackupFile := getPathToAwsCredentialsBackupFile()
	err := os.Remove(pathToCredentialsBackupFile)

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Warning: Unable to remove old backup of credentials file (%s): %s\n",
			pathToCredentialsBackupFile,
			err.Error(),
		)
	} else {
		fmt.Printf("Removed old backup of credentials file (%s)\n", pathToCredentialsBackupFile)
	}
}

func getPathToAwsDirectory() string {
	u, err := user.Current()

	if err != nil {
		exitWithErrorMessage("Unable to determine current user: %s\n", err.Error())
	}

	pathToHomeDirectory := u.HomeDir

	return path.Join(pathToHomeDirectory, ".aws")
}

func getPathToAwsCredentialsFile() string {
	return path.Join(getPathToAwsDirectory(), "credentials")
}

func getPathToAwsCredentialsBackupFile() string {
	return path.Join(getPathToAwsDirectory(), nameOfCredentialsBackupFile)
}

func doesCredentialsFileExist() bool {
	pathToCredentialsFile := getPathToAwsCredentialsFile()

	return doesFileExist(pathToCredentialsFile)
}

func doesFileExist(pathToFile string) bool {
	_, err := os.Stat(pathToFile)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		exitWithErrorMessage("Unable to check if file exists (%s): %s\n", pathToFile, err.Error())
	}

	return true
}

func doesCredentialsFileDefaultProfileContainPermanentCredentials() bool {
	pathToCredentialsFile := getPathToAwsCredentialsFile()
	credentialsFileContent, err := ioutil.ReadFile(pathToCredentialsFile)
	const defaultFailureMessageFormat = "Unable to determine if default profile in credentials file contains permanent credentials: %s\n"

	if err != nil {
		exitWithErrorMessage(defaultFailureMessageFormat, err.Error())
	}

	credentialsConfig, err := ini.Load(credentialsFileContent)

	if err != nil {
		exitWithErrorMessage(defaultFailureMessageFormat, err.Error())
	}

	defaultProfile := getDefaultProfileFromCredentialsIniConfiguration(credentialsConfig)

	return doesProfileContainPermanentCredentials(defaultProfile)
}

func getDefaultProfileFromCredentialsIniConfiguration(configuration *ini.File) *ini.Section {
	profiles := configuration.Sections()

	if len(profiles) < 1 {
		return nil
	}

	for _, profile := range profiles {
		if profile.Name() == "default" {
			return profile
		}
	}

	return nil
}

func doesProfileContainPermanentCredentials(profile *ini.Section) bool {
	accessKeyIDConfigurationKey, err := profile.GetKey("aws_access_key_id")

	if err != nil {
		return false
	}

	accessKeyID := accessKeyIDConfigurationKey.Value()

	if doesAccessKeyIDValueIndicatePermanentAccessKey(accessKeyID) && false == doesProfileContainSessionToken(profile) {
		return true
	}

	return false
}

func doesAccessKeyIDValueIndicatePermanentAccessKey(accessKeyID string) bool {
	return strings.HasPrefix(accessKeyID, "AKIA")
}

func doesProfileContainSessionToken(profile *ini.Section) bool {
	_, err := profile.GetKey("aws_session_token")

	if err != nil {
		return false
	}

	return true
}
