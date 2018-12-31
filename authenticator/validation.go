package authenticator

import (
	"errors"
	"strconv"
)

func ValidateMFATokenFormat(mfaToken string) error {
	_, err := strconv.ParseUint(mfaToken, 10, 64)
	if err != nil {
		return errors.New("MFA token should be a number")
	}

	if len(mfaToken) != 6 {
		return errors.New("MFA token should be six digits long")
	}

	return nil
}
