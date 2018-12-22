package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func requestNewTemporaryCredentials(mfaToken string, durationInSeconds int64) *sts.Credentials {
	session := createAwsSession()
	stsClient := sts.New(session)

	mfaSerialNumber := determineMfaSerialNumber(stsClient)
	result, err := getSessionToken(session, mfaSerialNumber, mfaToken, durationInSeconds)

	if err != nil {
		exitWithFormattedErrorMessage("Authentication failed: %s\n", err.Error())
	}

	return result.Credentials
}

func determineMfaSerialNumber(stsClient *sts.STS) string {
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		exitWithFormattedErrorMessage("Unable to get caller identity: %s\n", err.Error())
	}

	awsAccountNumber := *callerIdentity.Account
	userName := getUserNameFromCallerIdentity(callerIdentity)

	return assembleMfaSerialNumberFromComponents(awsAccountNumber, userName)
}

func getUserNameFromCallerIdentity(callerIdentity *sts.GetCallerIdentityOutput) string {
	const separator = "/"

	return strings.Split(*callerIdentity.Arn, separator)[1]
}

func isValidMfaTokenValue(mfaToken string) bool {
	_, err := strconv.ParseUint(mfaToken, 10, 64)

	if err != nil {
		return false
	}

	return true
}

func getSessionToken(session *session.Session, mfaSerialNumber, mfaToken string, durationInSeconds int64) (*sts.GetSessionTokenOutput, error) {
	stsClient := sts.New(session)

	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(durationInSeconds),
		SerialNumber:    aws.String(mfaSerialNumber),
		TokenCode:       aws.String(mfaToken),
	}

	return stsClient.GetSessionToken(input)
}

func assembleMfaSerialNumberFromComponents(awsAccountNumber, userName string) string {
	// Note: should be an ARN
	return fmt.Sprintf("arn:aws:iam::%s:mfa/%s", awsAccountNumber, userName)
}

func createAwsSession() *session.Session {
	return session.Must(session.NewSession())
}

func willEnvironmentVariablesPreemptUseOfCredentialsFile() bool {
	accessKeyIDEnvironmentVariableValue := os.Getenv("AWS_ACCESS_KEY_ID")

	return len(accessKeyIDEnvironmentVariableValue) != 0
}
