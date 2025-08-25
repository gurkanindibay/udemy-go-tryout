package utils


import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("supersecretkey")

func GenerateToken(email string, userId int64) (string, error) {
	// Implementation for generating JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
		"userId":   userId,
	})

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Implementation for validating JWT token
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})
}