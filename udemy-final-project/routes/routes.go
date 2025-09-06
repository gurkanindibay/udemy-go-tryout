package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/middlewares"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func SetupRoutes(server *gin.Engine) {
	// Swagger UI route
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.GET("/users/:id/registrations", getUserRegistrations)
	authenticated.DELETE("/events/:id/register", cancelRegistration)
}
