package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Hash the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswordAndHashedPassword(hashedPassword string, password string) bool {
	// returns true; if the hash matches.
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
