package main

import (
	"errors"
	"fmt"
	"github.com/luhring/awsmfa/mfa"
	"os"
)

const profileName = "default"

// todo: expose these with --version flag
var (
	version   = "No version provided"
	commit    = "No commit provided"
	buildTime = "No build timestamp provided"
)


func main() {
	numberOfArgumentsPassedIn := len(os.Args) - 1
	errUnexpectedArguments := errors.New("unexpected argument(s) passed in, type 'awsmfa --help' to see correct syntax")



	if numberOfArgumentsPassedIn == 0 {
		displayHelpText()
		os.Exit(0)
	}

	if numberOfArgumentsPassedIn == 1 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			help()
		}

		if os.Args[1] == "--restore" || os.Args[1] == "-r" {
			restore()
		}

		mfaToken := os.Args[1]
		authenticate(mfaToken)
	}

	exitWithError(errUnexpectedArguments)
}

func help() {
	displayHelpText()
	os.Exit(0)
}

func restore() {
	mfaController, err := mfa.New(profileName)

	if err != nil {
		// todo: do something (better) with this error
		exitWithError(fmt.Errorf("could not make mfa controller"))
	}

	err = mfaController.ExpireCredentials()
	if err != nil {
		// todo: do something (better) with this error
		exitWithError(fmt.Errorf("could not expire credentials"))
	}


	os.Exit(0)
}

func authenticate(mfaToken string) {


	mfaController, err := mfa.New(profileName)

	if err != nil {
		// todo: do something (better) with this error
		exitWithError(fmt.Errorf("could not make mfa controller"))
	}

	err = mfaController.RequestCredentials(mfaToken)
	if err != nil {
		// todo: do something (better) with this error
		exitWithError(fmt.Errorf("could not request credentials"))
	}

	// todo: add all printing to stdout/err here

	os.Exit(0)
}

func exitWithError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
