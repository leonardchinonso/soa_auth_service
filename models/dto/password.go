package dto

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

// Password is a custom type for managing passwords
type Password string

// Validate checks that a password meets the password requirements
func (p Password) Validate() error {
	const minPasswordLength = 6
	if isValidPassword(string(p), minPasswordLength) {
		return fmt.Errorf("invalid password")
	}
	return nil
}

// IsEqualValue compares the string value of a password to the input password
func (p Password) IsEqualValue(password Password) bool {
	return string(p) == string(password)
}

// Hash hashes a password using Go's bcrypt and a preset cost
func (p Password) Hash() (string, error) {
	const passwordHashCost = 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), passwordHashCost)
	return string(bytes), err
}

// IsEqualHash compares a password and a hash to see if they're equivalent
func (p Password) IsEqualHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	return err == nil
}

// isValidPassword checks that a password meets the following requirements:
// - is not less than minLength
// - has at least one number
// - has at least one special character
// - has at least one upper case character
func isValidPassword(password string, minLength int) bool {
	if len(password) < minLength {
		return false
	}

	var hasNum, hasSpecChar, hasUpperCase bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNum = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecChar = true
		default:
			return false
		}
	}

	return hasNum && hasSpecChar && hasUpperCase
}
