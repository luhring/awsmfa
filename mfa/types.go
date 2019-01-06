package mfa

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

type token string

type profile = ini.Section

type pair struct {
	key string
	value string
}

type Controller struct {
	CredentialsFile *awsCredentialsFile
	Profile         string
	stsClient       *sts.STS
}

type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Profile         string
}

type awsCredentialsFile struct {
	Filename      string
	Configuration *ini.File
}
