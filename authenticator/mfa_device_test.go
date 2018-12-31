package authenticator

import "testing"

func TestComputeARNForVirtualMFADevice(t *testing.T) {
	testCases := []struct {
		inputAwsAccountNumber string
		inputUserName         string
		expectedOutput        string
	}{
		{
			inputAwsAccountNumber: "123123123123",
			inputUserName:         "tony.stark",
			expectedOutput:        "arn:aws:iam::123123123123:mfa/tony.stark",
		},
		{
			inputAwsAccountNumber: "123412341234",
			inputUserName:         "tony.stark@starkindustries.com",
			expectedOutput:        "arn:aws:iam::123412341234:mfa/tony.stark@starkindustries.com",
		},
		{
			inputAwsAccountNumber: "123456123456",
			inputUserName:         "x-ray",
			expectedOutput:        "arn:aws:iam::123456123456:mfa/x-ray",
		},
	}

	for _, testCase := range testCases {
		output := computeARNForVirtualMFADevice(testCase.inputAwsAccountNumber, testCase.inputUserName)

		if output != testCase.expectedOutput {
			t.Errorf("expected output was '%s' but actual output was '%s'", testCase.expectedOutput, output)
		}
	}
}
