package crypto

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Decrypt(hashpass, pass string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(pass))
	if err != nil {
		return false, err
	}
	return true, nil
}
