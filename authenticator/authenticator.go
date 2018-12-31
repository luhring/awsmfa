package authenticator

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/luhring/awsmfa/credentials"
	"github.com/luhring/awsmfa/credentials_file"
	"github.com/luhring/awsmfa/environment"
	"github.com/luhring/awsmfa/file_coordinator"
	"os"
	"strconv"
	"strings"
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

	newCredentialsFile, err := credentials_file.NewFromCredentials(
		newCredentials,
		a.fileCoordinator.SelectedProfileName,
		a.fileCoordinator.Env.PathToCredentialsFile(),
	)

	err = newCredentialsFile.Save()
	if err != nil {
		return err
	}

	fmt.Print("\nAuthentication successful 👍\n\nSaved new session credentials to credentials file\n")

	if environment.WillEnvironmentVariablesPreemptUseOfCredentialsFile() {
		_, _ = fmt.Fprintf(os.Stderr, "\nWarning: Because you currently have the environment variable 'AWS_ACCESS_KEY_ID' set, most AWS CLI tools will use the credentials from your environment variables and not from your credentials file, which is where we just saved your new session credentials.\n\nYou might receive 'Access Denied' errors when performing actions that require MFA until you remove your AWS environment variables.\n")

		return nil
	}

	return nil
}

func ValidateMFATokenFormat(mfaToken string) error {
	_, err := strconv.ParseUint(mfaToken, 10, 64)
	if err != nil {
		return errors.New("MFA token should be a number")
	}

	if len(mfaToken) != 6 {
		return errors.New("MFA token should be six digits long")
	}

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

	return convertCredentialsFromStsCredentials(result.Credentials), nil
}

func convertCredentialsFromStsCredentials(input *sts.Credentials) *credentials.Credentials {
	return &credentials.Credentials{
		AccessKeyID:     *input.AccessKeyId,
		SecretAccessKey: *input.SecretAccessKey,
		SessionToken:    *input.SessionToken,
	}
}

func (a *Authenticator) computeMFADeviceSerialNumber() (string, error) {
	callerIdentity, err := a.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return "", nil
	}

	awsAccountNumber := *callerIdentity.Account
	userName := getUserNameFromCallerIdentity(callerIdentity)

	return computeARNForVirtualMFADevice(awsAccountNumber, userName), nil
}

func getUserNameFromCallerIdentity(callerIdentity *sts.GetCallerIdentityOutput) string {
	const separator = "/"
	return strings.Split(*callerIdentity.Arn, separator)[1]
}

func computeARNForVirtualMFADevice(awsAccountNumber, userName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:mfa/%s", awsAccountNumber, userName)
}

func (a *Authenticator) getSessionToken(mfaToken, serialNumber string, sessionDurationInSeconds int64) (*sts.GetSessionTokenOutput, error) {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(sessionDurationInSeconds),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(mfaToken),
	}

	return a.stsClient.GetSessionToken(input)
}
