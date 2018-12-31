package environment

import (
	"fmt"
	"os"
)

const NameOfVariableForAccessKeyID = "AWS_ACCESS_KEY_ID"

type Environment struct {
	homeDirectory string
}

func MustInit() *Environment {
	homeDirectory, err := getHomeDirectory()

	if err != nil {
		exitWithError(err)
	}

	return &Environment{
		homeDirectory: homeDirectory,
	}
}

func WillEnvironmentVariablesPreemptUseOfCredentialsFile() bool {
	accessKeyID := os.Getenv(NameOfVariableForAccessKeyID)

	return len(accessKeyID) != 0
}

func exitWithError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
