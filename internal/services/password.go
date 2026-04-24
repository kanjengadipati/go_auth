package services

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const PasswordHashCost = 12
const MinimumPasswordLength = 8

var ErrWeakPassword = errors.New("password must be at least 8 characters")

func ValidatePasswordStrength(password string) error {
	if utf8.RuneCountInString(password) < MinimumPasswordLength {
		return ErrWeakPassword
	}
	return nil
}

func HashPassword(password string) (string, error) {
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
