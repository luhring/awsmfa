package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const defaultSessionDurationInSeconds = 21600 // 6 hours

func main() {
	// # awsmfa s3 ls 414395
	// $ awsmfa 434134 s3 ls
	// $ awsmfa 321654

	// Were no arguments passed in?
	// Yes -- display help

	// Was there only one argument passed in?
	// Yes -- Was it a valid token?
	// 			Yes -- attempt to get session token and save as credentials

	//			No -- error out

	// No -- Was the first argument a valid token?
	//			Yes -- attempt to get session token. Success?
	//					Yes -- pass remaining arguments to new shell process spawned from "aws"
	//					No -- error out
	//			No -- error out

	numberOfArgumentsPassedIn := len(os.Args) - 1

	if numberOfArgumentsPassedIn == 0 {
		// display help

		os.Exit(0)
	}

	if numberOfArgumentsPassedIn == 1 {
		mfaToken := os.Args[1]

		if false == canMfaTokenBeConvertedToPositiveInteger(mfaToken) {
			fmt.Println("error: expected argument to be MFA token (integer)")
			os.Exit(1)
		}

		session := createAwsSession()
		stsClient := sts.New(session)

		callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

		if err != nil {
			fmt.Printf("unable to get caller identity: %s\n", err.Error())
			os.Exit(1)
		}

		awsAccountNumber := *callerIdentity.Account
		userName := getUserNameFromCallerIdentity(callerIdentity)
		mfaSerialNumber := computeMfaSerialNumber(awsAccountNumber, userName)

		result, err := getSessionToken(session, mfaSerialNumber, mfaToken, defaultSessionDurationInSeconds)

		if err != nil {
			fmt.Printf("unable to get session token: %s\n", err.Error())
			os.Exit(1)
		}

		newCredentials := result.Credentials
		credentialsFileContent := generateCredentialsFileContent(newCredentials)

		backupCredentialsFileAndSaveNewCredentialsToDisk(credentialsFileContent)

		fmt.Println("Successfully obtained session credentials.")
		os.Exit(0)
	}
}

func getUserNameFromCallerIdentity(callerIdentity *sts.GetCallerIdentityOutput) string {
	const separator = "/"

	return strings.Split(*callerIdentity.Arn, separator)[1]
}

func canMfaTokenBeConvertedToPositiveInteger(mfaToken string) bool {
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

func computeMfaSerialNumber(awsAccountNumber, userName string) string {
	// Note: should be an ARN
	// e.g.: arn:aws:iam::120123456789:mfa/josh.is.crazy

	return fmt.Sprintf("arn:aws:iam::%s:mfa/%s", awsAccountNumber, userName)
}

func createAwsSession() *session.Session {
	return session.Must(session.NewSession())
}

func generateCredentialsFileContent(credentials *sts.Credentials) string {
	return fmt.Sprintf(`[default]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s
`, *credentials.AccessKeyId, *credentials.SecretAccessKey, *credentials.SessionToken)
}

func backupCredentialsFileAndSaveNewCredentialsToDisk(newCredentialsFileContent string) {
	pathToCredentialsFile := path.Join(getPathToAwsDirectory(), "credentials")

	if _, err := os.Stat(pathToCredentialsFile); false == os.IsNotExist(err) {
		const nameOfBackupCredentialsFile = "credentials_backup_by_awsmfa"
		pathToBackupCredentialsFile := path.Join(getPathToAwsDirectory(), nameOfBackupCredentialsFile)

		err = os.Rename(pathToCredentialsFile, pathToBackupCredentialsFile)

		if err != nil {
			fmt.Printf("unable to back up AWS credentials file: %s", err.Error())
			os.Exit(1)
		}
	}

	err := ioutil.WriteFile(pathToCredentialsFile, []byte(newCredentialsFileContent), 0600)

	if err != nil {
		fmt.Printf("unable to save new sesion credentials to AWS credentials file: %s\n", err.Error())
		os.Exit(1)
	}
}

func getPathToAwsDirectory() string {
	u, err := user.Current()

	if err != nil {
		fmt.Printf("unable to determine current user: %s\n", err.Error())
		os.Exit(1)
	}

	pathToHomeDirectory := u.HomeDir

	return path.Join(pathToHomeDirectory, ".aws")
}
