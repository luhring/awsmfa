package authenticator

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/luhring/awsmfa/credentials"
	"github.com/luhring/awsmfa/credentials_file"
	"github.com/luhring/awsmfa/environment"
	"github.com/luhring/awsmfa/file_coordinator"
	"os"
)

const defaultSessionDurationInSeconds = 21600 // 6 hours

type Authenticator struct {
	stsClient       *sts.STS
	fileCoordinator *file_coordinator.Coordinator
}

func New(stsClient *sts.STS, fileCoordinator *file_coordinator.Coordinator) (*Authenticator, error) {
	return &Authenticator{
		stsClient,
		fileCoordinator,
	}, nil
}

func (a *Authenticator) AuthenticateUsingMFA(mfaToken string) error {
	newCredentials, err := a.requestNewTemporaryCredentials(mfaToken, defaultSessionDurationInSeconds)
	if err != nil {
		return err
	}

	fmt.Println("Multi-factor authentication was successful")

	newCredentialsFile, err := credentials_file.NewFromCredentials(
		newCredentials,
		a.fileCoordinator.SelectedProfileName,
		a.fileCoordinator.Env.PathToCredentialsFile(),
	)
	if err != nil {
		return err
	}

	err = newCredentialsFile.Save()
	if err != nil {
		return err
	}

	fmt.Println("Saved new session credentials to credentials file")

	if environment.WillEnvironmentVariablesPreemptUseOfCredentialsFile() {
		_, _ = fmt.Fprintf(os.Stderr, "\nWARNING: Because you have the environment variable '%s' set, most AWS tools will use the credentials from your environment variables and not from your credentials file, which is where we just saved your new session credentials.\n\nYou might receive 'Access Denied' errors when performing actions that require MFA until you remove your AWS environment variables.\n", environment.NameOfVariableForAccessKeyID)

		return nil
	}

	fmt.Print("\nYou now have access to actions where your IAM policies require 'MultiFactorAuthPresent' 👍\n")

	return nil
}

func (a *Authenticator) requestNewTemporaryCredentials(mfaToken string, sessionDurationInSeconds int64) (*credentials.Credentials, error) {
	serialNumber, err := a.computeMFADeviceSerialNumber()
	if err != nil {
		return nil, err
	}

	result, err := a.getSessionToken(mfaToken, serialNumber, sessionDurationInSeconds)

	if err != nil {
		return nil, err
	}

	return credentials.New(
		*result.Credentials.AccessKeyId,
		*result.Credentials.SecretAccessKey,
		*result.Credentials.SessionToken,
	), nil
}

func (a *Authenticator) getSessionToken(mfaToken, serialNumber string, sessionDurationInSeconds int64) (*sts.GetSessionTokenOutput, error) {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(sessionDurationInSeconds),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(mfaToken),
	}

	return a.stsClient.GetSessionToken(input)
}
