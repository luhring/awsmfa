package credentials

import (
	"strings"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func New(accessKeyID, secretAccessKey, sessionToken string) *Credentials {
	return &Credentials{
		accessKeyID,
		secretAccessKey,
		sessionToken,
	}
}

func (c *Credentials) ArePermanent() bool {
	return c.havePermanentAccessKeyID() && false == c.HaveSessionToken()
}

func (c *Credentials) HaveSessionToken() bool {
	return c.SessionToken != ""
}

func (c *Credentials) havePermanentAccessKeyID() bool {
	const accessKeyIDPrefix = "AKIA"

	return strings.HasPrefix(c.AccessKeyID, accessKeyIDPrefix)
}
