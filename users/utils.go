package users

import (
	"fmt"
	"regexp"
	"unicode"
)

const passwordLength = 8

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//noinspection GoErrorStringFormat
var (
	ErrPasswordLen     = fmt.Sprintf("Needs at least %d characters", passwordLength)
	ErrPasswordUpper   = "Needs at least 1 uppercase character"
	ErrPasswordLower   = "Needs at least 1 lowercase character"
	ErrPasswordNumber  = "Needs at least 1 number"
	ErrPasswordSpecial = "Needs at least 1 special character"
)

func ValidatePassword(password string) (bool, []string) {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	errors := make([]string, 0)
	if len(password) < passwordLength {
		errors = append(errors, ErrPasswordLen)
	}
	if !hasUpper {
		errors = append(errors, ErrPasswordUpper)
	}
	if !hasLower {
		errors = append(errors, ErrPasswordLower)
	}
	if !hasNumber {
		errors = append(errors, ErrPasswordNumber)
	}
	if !hasSpecial {
		errors = append(errors, ErrPasswordSpecial)
	}

	if len(errors) > 0 {
		return false, errors
	}
	return true, nil
}
