package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/models"
	"github.com/gurkanindibay/udemy-rest-api/utils"
	"log"
)


func registerUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := user.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func loginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//log email and password
	log.Printf("Attempting login for user: %s", user.Email)
	log.Printf("User password: %s", user.Password)

	verifiedUser, err := models.VerifyUserCredentials(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if verifiedUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(verifiedUser.Email, verifiedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}
