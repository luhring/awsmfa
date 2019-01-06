package mfa

import (
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
)

const (
	keyAccessKeyID     = "aws_access_key_id"
	keySecretAccessKey = "aws_secret_access_key"
	keySessionToken    = "aws_session_token"
	keyBackupPrefix    = "backup_"
)

func newAWSCredentialsFile(filename string) (*awsCredentialsFile, error) {
	// todo: if file doens't exist, create it and don't return an error
	// todo: should not be doing file IO in constructor

	credentialsFileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	configuration, err := ini.Load(credentialsFileContent)
	if err != nil {
		return nil, err
	}

	return &awsCredentialsFile{
		Filename:      filename,
		Configuration: configuration,
	}, nil
}

func (credFile *awsCredentialsFile) Save() error {
	return credFile.Configuration.SaveTo(credFile.Filename)
}

func (credFile *awsCredentialsFile) backup(profile string) error {
	currentBackupCredentials, err := credFile.getBackupCredentials(profile)
	if err == nil && currentBackupCredentials != nil {
		return fmt.Errorf("backup credentials already exist")
	}

	currentCredentials, err := credFile.getCredentials(profile)
	if err != nil || currentCredentials == nil {
		return fmt.Errorf("no permanent credentials to backup")
	}

	return credFile.setBackupCredentials(currentCredentials)
}

func (credFile *awsCredentialsFile) restore(profile string) error {
	currentBackupCredentials, err := credFile.getBackupCredentials(profile)
	if err != nil || currentBackupCredentials == nil {
		return fmt.Errorf("no backup credentials exist")
	}

	currentCredentials, err := credFile.getCredentials(profile)
	if err == nil && currentCredentials != nil {
		return fmt.Errorf("permanent credentials already exist")
	}

	return credFile.setCredentials(currentBackupCredentials)
}

func (credFile *awsCredentialsFile) setCredentials(creds *awsCredentials) error {
	pairs := []pair{
		{key: keyAccessKeyID, value: creds.AccessKeyID},
		{key: keySecretAccessKey, value: creds.SecretAccessKey},
		{key: keySessionToken, value: creds.SessionToken},
	}

	return credFile.setProfileValues(creds.Profile, pairs...)
}

func (credFile *awsCredentialsFile) setBackupCredentials(creds *awsCredentials) error {
	pairs := []pair{
		{key: keyBackupPrefix+keyAccessKeyID, value: creds.AccessKeyID},
		{key: keyBackupPrefix+keySecretAccessKey, value: creds.SecretAccessKey},
		{key: keyBackupPrefix+keySessionToken, value: creds.SessionToken},
	}

	return credFile.setProfileValues(creds.Profile, pairs...)
}

func (credFile *awsCredentialsFile) getCredentials(profile string) (*awsCredentials, error) {
	values := credFile.getProfileKeys(profile, keyAccessKeyID, keySecretAccessKey, keySessionToken)
	return newAWSCredentials(
		values[keyAccessKeyID],
		values[keySecretAccessKey],
		values[keySessionToken],
		profile,
	)
}

func (credFile *awsCredentialsFile) getBackupCredentials(profile string) (*awsCredentials, error) {
	values := credFile.getProfileKeys(profile, keyBackupPrefix+keyAccessKeyID, keyBackupPrefix+keySecretAccessKey, keyBackupPrefix+keySessionToken)

	return newAWSCredentials(
		values[keyBackupPrefix+keyAccessKeyID],
		values[keyBackupPrefix+keySecretAccessKey],
		values[keyBackupPrefix+keySessionToken],
		profile,
	)
}

func (credFile *awsCredentialsFile) getProfileKeys(profile string, keys... string) map[string]string {
	result := make(map[string]string, 0)
	p := credFile.getProfile(profile)
	if p == nil {
		// todo: should this return an error instead?
		return result
	}
	for _, key := range keys {
		result[key] = ""
		item, err := p.GetKey(key)
		if err == nil {
			result[key] = item.Value()
		}
	}
	return result
}

func (credFile *awsCredentialsFile) getProfile(name string) *profile {
	profiles := credFile.Configuration.Sections()

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

func (credFile *awsCredentialsFile) setProfileValues(profile string, pairs... pair) error {
	// todo: if no profile, use default

	p := credFile.getProfile(profile)
	if p == nil {
		return fmt.Errorf("cannot find profile")
	}

	for _, pair := range pairs {
		if pair.value == "" {
			continue
		}
		if p.HasKey(pair.key) {
			p.Key(pair.key).SetValue(pair.value)
		} else {
			_, err := p.NewKey(pair.key, pair.value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
