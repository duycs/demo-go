package helpers

import (
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GeneratePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
