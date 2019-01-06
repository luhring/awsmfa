package mfa

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"os/user"
	"path"
	"strings"
)

const profileName = "default"
const defaultSessionDurationInSeconds = 21600 // 6 hours
const envKeyAccessKeyId = "AWS_ACCESS_KEY_ID"

// todo: add a constructor that takes a credentials path override



// todo: add profile switch to application (or detect from environment)
func New(profile string) (*Controller, error) {
	path, err := pathToCredentialsFile()
	if err != nil {
		return nil, err
	}

	file, err := newAWSCredentialsFile(path)
	if err != nil {
		return nil, err
	}
	return &Controller{
		CredentialsFile: file,
		Profile: profile,
		stsClient: sts.New(session.Must(session.NewSession())),
	}, nil
}

// todo: make a constructor that can behave appropriately when this should be driven by AWS_ env vars
// func NewFromEnvironment() (*Controller, error) {
// 	&Controller{
// 		CredentialsFile: .... new,
// 		profile: .... from env or default,
// 		stsClient: sts.New(session.Must(session.NewSession())),
// 	}, nil
// }



func WillEnvironmentVariablesPreemptUseOfCredentialsFile() bool {
	accessKeyID := os.Getenv(envKeyAccessKeyId)

	return len(accessKeyID) != 0
}

func pathToCredentialsFile() (string, error) {
	u, err := user.Current()

	if err != nil {
		return "", err
	}

	return path.Join(u.HomeDir, ".aws", "credentials"), nil
}



// todo: probably should not be writing to screen from here (the caller should based on the error return)
func (controller *Controller) RequestCredentials(mfaToken string) error {
	tempCredentials, err := controller.requestCredentials(mfaToken, defaultSessionDurationInSeconds)
	if err != nil {
		return err
	}

	fmt.Println("Multi-factor authentication was successful")

	err = controller.CredentialsFile.backup(profileName)
	if err != nil {
		return err
	}

	err = controller.CredentialsFile.setCredentials(tempCredentials)
	if err != nil {
		return err
	}

	err = controller.CredentialsFile.Save()
	if err != nil {
		return err
	}

	fmt.Println("Saved new session credentials to credentials file")

	if WillEnvironmentVariablesPreemptUseOfCredentialsFile() {
		// todo: return error for the caller to use
		_, _ = fmt.Fprintf(os.Stderr, "\nWARNING: Because you have the environment variable '%s' set, most AWS tools will use the credentials from your environment variables and not from your credentials file, which is where we just saved your new session credentials.\n\nYou might receive 'Access Denied' errors when performing actions that require MFA until you remove your AWS environment variables.\n", envKeyAccessKeyId)

		return nil
	}

	fmt.Print("\nYou now have access to actions where your IAM policies require 'MultiFactorAuthPresent' 👍\n")

	return nil
}

func (controller *Controller) ExpireCredentials() error {
	err := controller.CredentialsFile.restore(profileName)
	if err != nil {
		return err
	}
	return controller.CredentialsFile.Save()
}



func (controller *Controller) requestCredentials(token string, sessionDurationInSeconds int64) (*awsCredentials, error) {
	mfaToken, err := newMFAToken(token)
	if err != nil {
		return nil, err
	}

	serialNumber, err := controller.mfaDeviceSerialNumber()
	if err != nil {
		return nil, err
	}

	result, err := controller.getSessionToken(mfaToken, serialNumber, sessionDurationInSeconds)

	if err != nil {
		return nil, err
	}

	return &awsCredentials{
		AccessKeyID: *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken: *result.Credentials.SessionToken,
		Profile: "",
	}, nil
}

func (controller *Controller) getSessionToken(token token, serialNumber string, sessionDurationInSeconds int64) (*sts.GetSessionTokenOutput, error) {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(sessionDurationInSeconds),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(string(token)),
	}

	return controller.stsClient.GetSessionToken(input)
}

func (controller *Controller) mfaDeviceSerialNumber() (string, error) {
	callerIdentity, err := controller.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return "", err
	}

	userName := strings.Split(*callerIdentity.Arn, "/")[1]
	serialArn := fmt.Sprintf("arn:aws:iam::%s:mfa/%s", *callerIdentity.Account, userName)

	return serialArn, nil
}


