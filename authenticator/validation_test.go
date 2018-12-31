package authenticator

import (
	"errors"
	"testing"
)

func TestValidateMFATokenFormat(t *testing.T) {
	testCases := []struct {
		mfaToken       string
		expectedOutput error
	}{
		{
			mfaToken:       "123456",
			expectedOutput: nil,
		},
		{
			mfaToken:       "1234567",
			expectedOutput: errors.New(ErrMFATokenIncorrectLength),
		},
		{
			mfaToken:       "1",
			expectedOutput: errors.New(ErrMFATokenIncorrectLength),
		},
		{
			mfaToken:       "hello",
			expectedOutput: errors.New(ErrMFATokenNotANumber),
		},
		{
			mfaToken:       "-1",
			expectedOutput: errors.New(ErrMFATokenNotANumber),
		},
		{
			mfaToken:       "--hello-there",
			expectedOutput: errors.New(ErrMFATokenNotANumber),
		},
	}

	for _, testCase := range testCases {
		output := ValidateMFATokenFormat(testCase.mfaToken)

		if output != testCase.expectedOutput {
			if testCase.expectedOutput == nil {
				t.Errorf("expected no error but received error: %v -- mfaToken was '%s'", output, testCase.mfaToken)
			} else if output == nil {
				t.Errorf("received no error but expected error: %v -- mfaToken was '%s'", testCase.expectedOutput, testCase.mfaToken)
			} else if output.Error() != testCase.expectedOutput.Error() {
				t.Errorf("expected error '%v' but received error '%v' -- mfaToken was '%s'", testCase.expectedOutput, output, testCase.mfaToken)
			}
		}
	}
}
