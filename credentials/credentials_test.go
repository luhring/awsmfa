package credentials

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	accessKeyID := "AKIA-something-something"
	secretAccessKey := "something-something-something"
	sessionToken := "token-of-the-session"
	timeNow := time.Now()

	testStsCredentials := &sts.Credentials{
		AccessKeyId:     &accessKeyID,
		SecretAccessKey: &secretAccessKey,
		SessionToken:    &sessionToken,
		Expiration:      &timeNow,
	}

	c := New(
		*testStsCredentials.AccessKeyId,
		*testStsCredentials.SecretAccessKey,
		*testStsCredentials.SessionToken,
	)

	if c.AccessKeyID != accessKeyID {
		t.Error("AccessKeyID was not set correctly")
	}

	if c.SecretAccessKey != secretAccessKey {
		t.Error("SecretAccessKey was not set correctly")
	}

	if c.SessionToken != sessionToken {
		t.Error("SessionToken was not set correctly")
	}
}

func TestHasPermanentAccessKeyID(t *testing.T) {
	testCases := []struct {
		c              *Credentials
		expectedOutput bool
	}{
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			true,
		},
		{
			&Credentials{
				"ASIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			false,
		},
		{
			&Credentials{
				"ABABCBABSBSBSABCAKIA",
				"secret-access-key",
				"",
			},
			false,
		},
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"session-token",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		output := testCase.c.hasPermanentAccessKeyID()

		if output != testCase.expectedOutput {
			t.Errorf("evaluating %v; expected %v but got %v", testCase.c, testCase.expectedOutput, output)
		}
	}
}

func TestHasSessionToken(t *testing.T) {
	testCases := []struct {
		c              *Credentials
		expectedOutput bool
	}{
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			false,
		},
		{
			&Credentials{
				"ASIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			false,
		},
		{
			&Credentials{
				"ASIAABABCBABSBSBSABC",
				"secret-access-key",
				"session-token",
			},
			true,
		},
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"session-token",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		output := testCase.c.HasSessionToken()

		if output != testCase.expectedOutput {
			t.Errorf("evaluating %v; expected %v but got %v", testCase.c, testCase.expectedOutput, output)
		}
	}
}

func TestArePermanent(t *testing.T) {
	testCases := []struct {
		c              *Credentials
		expectedOutput bool
	}{
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			true,
		},
		{
			&Credentials{
				"ASIAABABCBABSBSBSABC",
				"secret-access-key",
				"",
			},
			false,
		},
		{
			&Credentials{
				"ASIAABABCBABSBSBSABC",
				"secret-access-key",
				"session-token",
			},
			false,
		},
		{
			&Credentials{
				"AKIAABABCBABSBSBSABC",
				"secret-access-key",
				"session-token",
			},
			false,
		},
	}

	for _, testCase := range testCases {
		output := testCase.c.ArePermanent()

		if output != testCase.expectedOutput {
			t.Errorf("evaluating %v; expected %v but got %v", testCase.c, testCase.expectedOutput, output)
		}
	}
}
