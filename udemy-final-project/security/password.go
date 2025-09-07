// Package security provides functions related to password hashing and
// verification used across the application authentication flows.
package security

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash for the given password.
// A cost of 14 is used which may be adjusted in the future for performance.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verifies if a password matches its bcrypt hash.
// It returns an error when the hash does not match or if the hash is invalid.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
