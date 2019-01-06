package mfa

import (
	"errors"
	"strconv"
)

const (
	ErrMFATokenNotANumber      = "MFA token should be a number"
	ErrMFATokenIncorrectLength = "MFA token should be six digits long"
)

func newMFAToken(tokenStr string) (token, error) {
	mfaToken := token(tokenStr)
	return mfaToken, mfaToken.Validate()
}

func (token token) Validate() error {
	_, err := strconv.ParseUint(string(token), 10, 64)
	if err != nil {
		return errors.New(ErrMFATokenNotANumber)
	}

	if len(string(token)) != 6 {
		return errors.New(ErrMFATokenIncorrectLength)
	}

	return nil
}
