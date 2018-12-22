package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const defaultSessionDurationInSeconds = 21600 // 6 hours
const nameOfCredentialsBackupFile = "credentials_backup_by_awsmfa"

var (
	version   = "No version provided"
	commit    = "No commit provided"
	buildTime = "No build timestamp provided"
)

func main() {
	numberOfArgumentsPassedIn := len(os.Args) - 1

	if numberOfArgumentsPassedIn == 0 {
		displayHelpText()
	}

	if numberOfArgumentsPassedIn == 1 {
		if os.Args[1] == "--restore" || os.Args[1] == "-r" {
			restorePermanentCredentials()
			os.Exit(0)
		}

		handlePersistentAuthenticationProcess()
		os.Exit(0)
	}
}

func exitWithFormattedErrorMessage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format, a...)
	os.Exit(1)
}

func displayHelpText() {
	const helpText = `To obtain temporary session credentials from AWS and save to credentials file:
awsmfa <mfa-token>

To restore original AWS credentials to credentials file:
awsmfa --restore
`
	fmt.Println(helpText)
	os.Exit(0)
}

func restorePermanentCredentials() {
	if doesCredentialsFileExist() {
		if doesCredentialsFileDefaultProfileContainPermanentCredentials() {
			fmt.Printf("'default' profile in %s already contains original credentials.\n", getPathToAwsCredentialsFile())
			removeCredentialsBackupFileIfItExists()

			return
		}

		if doesCredentialsBackupFileExist() {
			restoreCredentialsFileFromBackup()

			return
		}

		exitWithFormattedErrorMessage("Unable to find original (non-temporary) credentials!\n")
	}

	if doesCredentialsBackupFileExist() {
		restoreCredentialsFileFromBackup()

		return
	}

	exitWithFormattedErrorMessage("Unable to find any AWS credentials.\n")
}

func handlePersistentAuthenticationProcess() {
	mfaToken := os.Args[1]

	if false == isValidMfaTokenValue(mfaToken) {
		exitWithFormattedErrorMessage("Expected argument to be MFA token (integer).\n")
	}

	prepareCredentialsFileForUse()

	newCredentials := requestNewTemporaryCredentials(mfaToken, defaultSessionDurationInSeconds)
	newCredentialsFileContent := generateCredentialsFileContent(newCredentials)

	if doesCredentialsFileExist() {
		backUpCredentialsFile()
	}

	pathToCredentialsFile := getPathToAwsCredentialsFile()
	err := ioutil.WriteFile(pathToCredentialsFile, []byte(newCredentialsFileContent), 0600)

	if err != nil {
		exitWithFormattedErrorMessage("Unable to save new session credentials to %s: %s\n", pathToCredentialsFile, err.Error())
	}

	fmt.Printf("Authentication successful! Saved new session credentials to %s.\n", pathToCredentialsFile)

	if willEnvironmentVariablesPreemptUseOfCredentialsFile() {
		fmt.Fprintf(os.Stderr, "\nWarning: Because you currently have the environment variable 'AWS_ACCESS_KEY_ID' set, most AWS CLI tools will use the credentials from your environment variables and not the session credentials you just received, which are saved at %s.\n\n", pathToCredentialsFile)
	}
}
