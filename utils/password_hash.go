package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateHash accepts password, hash it and return it with error
func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword accepts password and hashed string and return boolean of wether they are equal
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
