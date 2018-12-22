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
		handlePersistentAuthenticationAttempt()
	}
}

func exitWithErrorMessage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format, a...)
	os.Exit(1)
}

func displayHelpText() {
	const helpText = `usage: awsmfa <mfa-token>`
	fmt.Println(helpText)
	os.Exit(0)
}

func handlePersistentAuthenticationAttempt() {
	mfaToken := os.Args[1]

	if false == isValidMfaTokenValue(mfaToken) {
		exitWithErrorMessage("Expected argument to be MFA token (integer)\n")
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
		exitWithErrorMessage("Unable to save new session credentials to %s: %s\n", pathToCredentialsFile, err.Error())
	}

	fmt.Printf("Authentication successful! Saved new session credentials to %s\n", pathToCredentialsFile)
	os.Exit(0)
}
