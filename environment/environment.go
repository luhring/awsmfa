package environment

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

const (
	NameOfVariableForAccessKeyID = "AWS_ACCESS_KEY_ID"
	nameOfCredentialsFile        = "credentials"
	nameOfCredentialsFileBackup  = "credentials_backup_by_awsmfa"
	nameOfAwsDirectory           = ".aws"
)

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

func (e *Environment) DoesHaveCredentialsFile() bool {
	return doesFileExist(e.PathToCredentialsFile())
}

func (e *Environment) DoesHaveCredentialsFileBackup() bool {
	return doesFileExist(e.PathToCredentialsFileBackup())
}

func WillEnvironmentVariablesPreemptUseOfCredentialsFile() bool {
	accessKeyID := os.Getenv(NameOfVariableForAccessKeyID)

	return len(accessKeyID) != 0
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

func (e *Environment) PathToCredentialsFile() string {
	return path.Join(e.pathToAwsDir(), nameOfCredentialsFile)
}

func (e *Environment) PathToCredentialsFileBackup() string {
	return path.Join(e.pathToAwsDir(), nameOfCredentialsFileBackup)
}

func (e *Environment) pathToAwsDir() string {
	return path.Join(e.homeDirectory, nameOfAwsDirectory)
}

func getHomeDirectory() (string, error) {
	u, err := user.Current()

	if err != nil {
		return "", err
	}

	return u.HomeDir, nil
}

func exitWithError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
