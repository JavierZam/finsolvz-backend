package utils

import (
	"crypto/rand"
	"encoding/hex"

	"finsolvz-backend/internal/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("PASSWORD_HASH_ERROR", "Failed to hash password", 500, err, nil)
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with plain text password
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("PASSWORD_MISMATCH", "Password does not match", 401, err, nil)
	}
	return nil
}

// GenerateRandomPassword generates a random password for forgot password functionality
func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 6) // 6 bytes = 12 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New("RANDOM_GENERATION_ERROR", "Failed to generate random password", 500, err, nil)
	}
	return hex.EncodeToString(bytes), nil
}
