package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
)

func TestGetUserNameFromCallerIdentity(t *testing.T) {
	testCases := []struct {
		arn            string
		expectedResult string
	}{
		{
			"arn:aws:iam::123123123123:user/tony",
			"tony",
		},
		{
			"arn:aws:iam::123123123123:user/tony.stark",
			"tony.stark",
		},
		{
			"arn:aws:iam::123123123123:user/tony@starkindustries.com",
			"tony@starkindustries.com",
		},
	}

	for _, testCase := range testCases {
		callerIdentity := &sts.GetCallerIdentityOutput{
			Arn: &testCase.arn,
		}

		result := getUserNameFromCallerIdentity(callerIdentity)

		if result != testCase.expectedResult {
			t.Errorf(
				"expected '%v' for '%v' but got '%v'",
				testCase.expectedResult,
				testCase.arn,
				result,
			)
		}
	}
}

func TestIsValidMfaTokenValue(t *testing.T) {
	testCases := []struct {
		mfaToken       string
		expectedResult bool
	}{
		{
			"123456",
			true,
		},
		{
			"4",
			true,
		},
		{
			"abcdef",
			false,
		},
		{
			"-12345",
			false,
		},
	}

	for _, testCase := range testCases {
		result := isValidMfaTokenValue(testCase.mfaToken)

		if result != testCase.expectedResult {
			t.Errorf(
				"expected '%v' for '%v' but got '%v'",
				testCase.expectedResult,
				testCase.mfaToken,
				result,
			)
		}
	}
}

func TestAssembleMfaSerialNumberFromComponents(t *testing.T) {
	testCases := []struct {
		awsAccountNumber string
		userName         string
		expectedResult   string
	}{
		{
			"12123123123",
			"tony",
			"arn:aws:iam::12123123123:mfa/tony",
		},
		{
			"555555555555",
			"mfa",
			"arn:aws:iam::555555555555:mfa/mfa",
		},
		{
			"123456123456",
			"h-e-l-l-o",
			"arn:aws:iam::123456123456:mfa/h-e-l-l-o",
		},
	}

	for _, testCase := range testCases {
		result := assembleMfaSerialNumberFromComponents(testCase.awsAccountNumber, testCase.userName)

		if result != testCase.expectedResult {
			t.Errorf(
				"expected '%v' but got '%v'",
				testCase.expectedResult,
				result,
			)
		}
	}
}
