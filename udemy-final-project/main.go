package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/routes"
	_ "github.com/gurkanindibay/udemy-rest-api/docs" // This is required for swagger
)

func main() {
	db.InitDB("events.db")
	server := gin.Default()
	
	// Add CORS middleware for Swagger UI
	server.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	
	routes.SetupRoutes(server)
	server.Run(":8080")
}

// Event Management API
//
// This is a REST API for managing events, user authentication, and event registrations.
//
// Terms Of Service: http://swagger.io/terms/
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions:
// BearerAuth:
//   type: apiKey
//   name: Authorization
//   in: header
//
// swagger:meta

