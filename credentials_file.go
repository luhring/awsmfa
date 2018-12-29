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

const nameOfCredentialsBackupFile = "credentials_backup_by_awsmfa"

func prepareCredentialsFileForUse() {
	if doesCredentialsFileExist() {
		if doesCredentialsFileDefaultProfileContainPermanentCredentials() { // as opposed to temporary credentials
			removeCredentialsBackupFileIfItExists()

			return
		}

		// Credentials file default profile has only temporary credentials.
		// We won't be able to use those credentials to get new temporary credentials.
		// Let's see if there's a backup file we can restore to the main credentials file.

		if doesCredentialsBackupFileExist() {
			restoreCredentialsFileFromBackup()
		}

		return
	}

	// We don't have usable credentials in the AWS credentials file.
	// Let's see if there's a backup file we can restore to the main credentials file.

	if doesCredentialsBackupFileExist() {
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

func backUpCredentialsFile() {
	credentialsFileBytes, err := ioutil.ReadFile(getPathToAwsCredentialsFile())

	if err != nil {
		errBackupFailed := fmt.Errorf("unable to back up credentials file (%s): %s", getPathToAwsCredentialsFile(), err.Error())
		exitWithError(errBackupFailed)
	}

	err = ioutil.WriteFile(getPathToAwsCredentialsBackupFile(), credentialsFileBytes, 0600)

	if err != nil {
		errBackupFailed := fmt.Errorf("unable to back up credentials file (%s): %s", getPathToAwsCredentialsFile(), err.Error())
		exitWithError(errBackupFailed)
	}

	fmt.Printf("Created backup of original credentials at %s.\n", getPathToAwsCredentialsBackupFile())
}

func restoreCredentialsFileFromBackup() {
	if doesCredentialsFileExist() {
		err := os.Remove(getPathToAwsCredentialsFile())

		if err != nil {
			errCredentialsRestoreFailed := fmt.Errorf("unable to restore AWS credentials file from backup: %s", err.Error())
			exitWithError(errCredentialsRestoreFailed)
		}
	}

	err := os.Rename(getPathToAwsCredentialsBackupFile(), getPathToAwsCredentialsFile())

	if err != nil {
		errCredentialsRestoreFailed := fmt.Errorf("unable to restore AWS credentials file from backup: %s", err.Error())
		exitWithError(errCredentialsRestoreFailed)
	}

	fmt.Printf("Restored original credentials from backup.\n")
}

func removeCredentialsBackupFile() {
	pathToCredentialsBackupFile := getPathToAwsCredentialsBackupFile()
	err := os.Remove(pathToCredentialsBackupFile)

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"warning: Unable to remove old backup of credentials file (%s): %s\n",
			pathToCredentialsBackupFile,
			err.Error(),
		)

		return
	}

	fmt.Printf("Deleted old backup of credentials file.\n")
}

func removeCredentialsBackupFileIfItExists() {
	if doesCredentialsBackupFileExist() {
		removeCredentialsBackupFile()
	}
}

func getPathToAwsDirectory() string {
	u, err := user.Current()

	if err != nil {
		errCannotDetermineCurrentUser := fmt.Errorf("unable to determine current user: %s", err.Error())
		exitWithError(errCannotDetermineCurrentUser)
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

func doesCredentialsBackupFileExist() bool {
	pathToCredentialsBackupFile := getPathToAwsCredentialsBackupFile()

	return doesFileExist(pathToCredentialsBackupFile)
}

func doesFileExist(pathToFile string) bool {
	_, err := os.Stat(pathToFile)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		errFileExistenceCheckFailed := fmt.Errorf("unable to check if file exists (%s): %s", pathToFile, err.Error())
		exitWithError(errFileExistenceCheckFailed)
	}

	return true
}

func doesCredentialsFileDefaultProfileContainPermanentCredentials() bool {
	pathToCredentialsFile := getPathToAwsCredentialsFile()
	credentialsFileContent, err := ioutil.ReadFile(pathToCredentialsFile)

	if err != nil {
		errProfileCheckForPermanentCredentialsFailed := fmt.Errorf("unable to determine if default profile in credentials file contains permanent credentials: %s", err.Error())
		exitWithError(errProfileCheckForPermanentCredentialsFailed)
	}

	credentialsConfig, err := ini.Load(credentialsFileContent)

	if err != nil {
		errProfileCheckForPermanentCredentialsFailed := fmt.Errorf("unable to determine if default profile in credentials file contains permanent credentials: %s", err.Error())
		exitWithError(errProfileCheckForPermanentCredentialsFailed)
	}

	defaultProfile := getDefaultProfileFromCredentialsIniConfiguration(credentialsConfig)

	return doesProfileContainPermanentCredentials(defaultProfile)
}

type Profile = ini.Section

func getDefaultProfileFromCredentialsIniConfiguration(configuration *ini.File) *Profile {
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
