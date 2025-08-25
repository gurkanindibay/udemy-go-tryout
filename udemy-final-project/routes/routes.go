package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/middlewares"
)

func SetupRoutes(server *gin.Engine) {
	// Public routes (no authentication required)
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEventByID)
	server.POST("/auth/register", registerUser)
	server.POST("/auth/login", loginUser)

	// Protected routes (authentication required)
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
}
