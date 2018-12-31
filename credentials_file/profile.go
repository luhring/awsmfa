package credentials_file

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/luhring/awsmfa/credentials"
)

const errFormatGettingCredentialsFromProfile = "unable to load AWS credentials from profile '%s'"

type Profile = ini.Section

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

func createNewKeyInProfile(profile *ini.Section, keyName, keyValue string) error {
	_, err := profile.NewKey(keyName, keyValue)

	if err != nil {
		return err
	}

	return nil
}
