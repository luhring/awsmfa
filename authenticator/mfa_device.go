package authenticator

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"strings"
)

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
