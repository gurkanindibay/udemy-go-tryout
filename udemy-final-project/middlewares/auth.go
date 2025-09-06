package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/utils"
)

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

	userId, err := utils.ValidateToken(tokenString)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	context.Set("userId", userId)
	context.Next()

}
