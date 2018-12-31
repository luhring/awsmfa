package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/luhring/awsmfa/authenticator"
	"github.com/luhring/awsmfa/environment"
	"github.com/luhring/awsmfa/file_coordinator"
	"os"
)

var (
	version   = "No version provided"
	commit    = "No commit provided"
	buildTime = "No build timestamp provided"
)

func main() {
	numberOfArgumentsPassedIn := len(os.Args) - 1
	errUnexpectedArguments := errors.New("unexpected argument(s) passed in, type 'awsmfa --help' to see correct syntax")

	const profileName = "default"

	if numberOfArgumentsPassedIn == 0 {
		displayHelpText()
		os.Exit(0)
	}

	if numberOfArgumentsPassedIn == 1 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			help()
		}

		env := environment.MustInit()
		fileCoordinator, err := file_coordinator.New(env, profileName)
		if err != nil {
			exitWithError(err)
		}

		if os.Args[1] == "--restore" || os.Args[1] == "-r" {
			restore(fileCoordinator)
		}

		mfaToken := os.Args[1]
		authenticate(fileCoordinator, mfaToken)
	}

	exitWithError(errUnexpectedArguments)
}

func authenticate(fileCoordinator *file_coordinator.Coordinator, mfaToken string) {
	err := authenticator.ValidateMFATokenFormat(mfaToken)
	if err != nil {
		exitWithError(err)
	}

	fileCoordinator.RestorePermanentCredentialsIfAppropriate()

	err = fileCoordinator.BackUpPermanentCredentialsIfPresent()
	if err != nil {
		exitWithError(err)
	}

	awsSession := session.Must(session.NewSession())
	stsClient := sts.New(awsSession)

	auth, err := authenticator.New(stsClient, fileCoordinator)
	if err != nil {
		exitWithError(err)
	}

	err = auth.AuthenticateUsingMFA(mfaToken)
	if err != nil {
		exitWithError(err)
	}

	os.Exit(0)
}

func help() {
	displayHelpText()
	os.Exit(0)
}

func restore(fileCoordinator *file_coordinator.Coordinator) {
	err := fileCoordinator.Restore()
	if err != nil {
		exitWithError(err)
	}

	os.Exit(0)
}

func exitWithError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
