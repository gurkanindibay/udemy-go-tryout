// Package security provides JWT generation and validation helpers.
package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NOTE: In production this key must come from a secure source (env/secret manager).
var jwtKey = []byte("supersecretkey")

// GenerateToken creates a JWT token for the given email and user ID.
func GenerateToken(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
		"userId": userID,
	})
	return token.SignedString(jwtKey)
}

// ValidateToken verifies a JWT token and returns the user ID if valid.
func ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}
	return int64(token.Claims.(jwt.MapClaims)["userId"].(float64)), nil
}
