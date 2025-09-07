// Package middlewares provides HTTP middleware functions for authentication and other common functionality.
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/security"
)

// Authenticate is a middleware that validates JWT tokens and sets user ID in context
func Authenticate(context *gin.Context) {
	// validate JWT token
	tokenString := context.GetHeader("Authorization")
	if tokenString == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := security.ValidateToken(tokenString)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	context.Set("userId", userID)
	context.Next()

}
