package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/middlewares"
	"github.com/gurkanindibay/udemy-rest-api/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	userService  services.UserService
	eventService services.EventService
	authService  services.AuthService
)

// InitServices initializes the service dependencies for the routes
func InitServices(u services.UserService, e services.EventService, a services.AuthService) {
	userService = u
	eventService = e
	authService = a
}

// SetupRoutes configures all the API routes for the application
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
