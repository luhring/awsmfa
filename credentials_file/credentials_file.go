package credentials_file

import (
	"errors"
	"github.com/go-ini/ini"
	"github.com/luhring/awsmfa/credentials"
	"io/ioutil"
	"os"
)

const (
	keyNameForAccessKeyID     = "aws_access_key_id"
	keyNameForSecretAccessKey = "aws_secret_access_key"
	keyNameForSessionToken    = "aws_session_token"
)

type CredentialsFile struct {
	Filename      string
	Configuration *ini.File
}

func NewFromDisk(filename string) (*CredentialsFile, error) {
	credentialsFileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	configuration, err := ini.Load(credentialsFileContent)

	if err != nil {
		return nil, err
	}

	return &CredentialsFile{
		Filename:      filename,
		Configuration: configuration,
	}, nil
}

func NewFromCredentials(c *credentials.Credentials, profileName, filename string) (*CredentialsFile, error) {
	configuration := ini.Empty()

	profile, err := configuration.NewSection(profileName)
	if err != nil {
		return nil, err
	}

	err = createNewKeyInProfile(profile, keyNameForAccessKeyID, c.AccessKeyID)
	if err != nil {
		return nil, err
	}

	err = createNewKeyInProfile(profile, keyNameForSecretAccessKey, c.SecretAccessKey)
	if err != nil {
		return nil, err
	}

	if c.HasSessionToken() {
		err = createNewKeyInProfile(profile, keyNameForSessionToken, c.SessionToken)

		if err != nil {
			return nil, err
		}
	}

	return &CredentialsFile{
		Filename:      filename,
		Configuration: configuration,
	}, nil
}

func NewFromConfiguration(configuration *ini.File, filename string) (*CredentialsFile, error) {
	if len(filename) == 0 {
		return nil, errors.New("filename parameter cannot be an empty string")
	}

	return &CredentialsFile{
		Filename:      filename,
		Configuration: configuration,
	}, nil
}

func (f *CredentialsFile) Save() error {
	return f.Configuration.SaveTo(f.Filename)
}

func (f *CredentialsFile) Delete() error {
	return os.Remove(f.Filename)
}
