package main

import (
	"fmt"
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
		// display help

		os.Exit(0)
	}

	if numberOfArgumentsPassedIn == 1 {
		mfaToken := os.Args[1]

		if false == isValidMfaTokenValue(mfaToken) {
			exitWithErrorMessage("Expected argument to be MFA token (integer)\n")
		}

		prepareCredentialsFileForUse()

		newCredentials := requestNewTemporaryCredentials(mfaToken, defaultSessionDurationInSeconds)
		newCredentialsFileContent := generateCredentialsFileContent(newCredentials)

		backupCredentialsFileAndSaveNewCredentialsToDisk(newCredentialsFileContent)

		os.Exit(0)
	}
}

func exitWithErrorMessage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format, a...)
	os.Exit(1)
}
