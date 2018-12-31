package authenticator

import (
	"errors"
	"strconv"
)

const (
	ErrMFATokenNotANumber      = "MFA token should be a number"
	ErrMFATokenIncorrectLength = "MFA token should be six digits long"
)

func ValidateMFATokenFormat(mfaToken string) error {
	_, err := strconv.ParseUint(mfaToken, 10, 64)
	if err != nil {
		return errors.New(ErrMFATokenNotANumber)
	}

	if len(mfaToken) != 6 {
		return errors.New(ErrMFATokenIncorrectLength)
	}

	return nil
}
