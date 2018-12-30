package credentials_file

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/luhring/awsmfa/credentials"
	"io/ioutil"
	"os"
)

const (
	keyNameForAccessKeyID                  = "aws_access_key_id"
	keyNameForSecretAccessKey              = "aws_secret_access_key"
	keyNameForSessionToken                 = "aws_session_token"
	errFormatGettingCredentialsFromProfile = "unable to load AWS credentials from profile '%s'"
)

type Profile = ini.Section

type CredentialsFile struct {
	Filename      string
	Configuration *ini.File
}

func LoadFromDisk(filename string) (*CredentialsFile, error) {
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

	if c.HaveSessionToken() {
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

func createNewKeyInProfile(profile *ini.Section, keyName, keyValue string) error {
	_, err := profile.NewKey(keyName, keyValue)

	if err != nil {
		return err
	}

	return nil
}

func (f *CredentialsFile) DoesProfileHavePermanentCredentials(profileName string) bool {
	c, err := f.GetCredentialsFromProfile(profileName)

	if err != nil {
		return false
	}

	return c.ArePermanent()
}

func (f *CredentialsFile) GetCredentialsFromProfile(name string) (*credentials.Credentials, error) {
	p := f.getProfile(name)

	accessKeyIDItem, err := p.GetKey(keyNameForAccessKeyID)
	accessKeyID := accessKeyIDItem.Value()

	if err != nil {
		return nil, fmt.Errorf(errFormatGettingCredentialsFromProfile, p.Name())
	}

	secretAccessKeyItem, err := p.GetKey(keyNameForSecretAccessKey)
	secretAccessKey := secretAccessKeyItem.Value()

	if err != nil {
		return nil, fmt.Errorf(errFormatGettingCredentialsFromProfile, p.Name())
	}

	sessionTokenItem, err := p.GetKey(keyNameForSessionToken)
	sessionToken := ""

	if err == nil {
		sessionToken = sessionTokenItem.Value()
	}

	return credentials.New(
		accessKeyID,
		secretAccessKey,
		sessionToken,
	), nil
}

func (f *CredentialsFile) getProfile(name string) *Profile {
	profiles := f.Configuration.Sections()

	if len(profiles) < 1 {
		return nil
	}

	for _, profile := range profiles {
		if profile.Name() == name {
			return profile
		}
	}

	return nil
}

func (f *CredentialsFile) Delete() error {
	return os.Remove(f.Filename)
}

func (f *CredentialsFile) LoadContentFrom(otherFile *CredentialsFile) {
	f.Configuration = otherFile.Configuration
}
