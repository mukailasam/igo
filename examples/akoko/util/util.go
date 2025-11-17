package util

import (
	"fmt"
	"regexp"
	"unicode"
)

const (
	userNameMinLength = 6
	userNameMaxLength = 20
)

func IsEmpty(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("all field required")
	}

	return nil
}

func ValidateUsername(username string) error {
	if len(username) < userNameMinLength || len(username) > userNameMaxLength {
		return fmt.Errorf("username must 6 characters minimun and 20 characters maximum")
	}

	return nil
}

func ValidateEmail(email string) error {

	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return fmt.Errorf("invalid email address")
	}

	return nil

}

func ValidatePassword(password string) error {
	var (
		hasMin     = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMin = true
	}

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else if unicode.IsPunct(char) {
			hasSpecial = true
		}
	}

	if !hasMin || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("paword should be minimum 8 in length and Password should contain at least a single uppercase letter, lowercase letter, single digit and a special character")
	}

	return nil
}
