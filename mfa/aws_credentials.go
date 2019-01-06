package mfa

import (
	"fmt"
	"strings"
)

// this indicates a permanent credential
const accessKeyIDPrefix = "AKIA"

// todo: handle case for no profile (default to "default")
func newAWSCredentials(accessKeyID, secretAccessKey, sessionToken, profile string) (*awsCredentials, error) {
	creds := &awsCredentials{
		accessKeyID,
		secretAccessKey,
		sessionToken,
		profile,
	}
	if !creds.IsValid() {
		return nil, fmt.Errorf("invalid credentials")
	}
	return creds, nil
}

func (cred *awsCredentials) ArePermanent() bool {
	return strings.HasPrefix(cred.AccessKeyID, accessKeyIDPrefix) && false == cred.HasSessionToken()
}

func (cred *awsCredentials) HasSessionToken() bool {
	return cred.SessionToken != ""
}

func (cred *awsCredentials) HasProfile() bool {
	return cred.Profile != ""
}

func (cred *awsCredentials) IsValid() bool {
	return cred.AccessKeyID != "" && cred.SecretAccessKey != ""
}
