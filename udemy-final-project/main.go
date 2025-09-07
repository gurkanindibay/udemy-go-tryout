// Package main provides the entry point for the event management API server.
package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/di"
	_ "github.com/gurkanindibay/udemy-rest-api/docs" // This is required for swagger
	"github.com/gurkanindibay/udemy-rest-api/grpc/auth"
	"github.com/gurkanindibay/udemy-rest-api/grpc/event"
	"github.com/gurkanindibay/udemy-rest-api/kafka"
	authpb "github.com/gurkanindibay/udemy-rest-api/proto/auth"
	eventpb "github.com/gurkanindibay/udemy-rest-api/proto/event"
	"github.com/gurkanindibay/udemy-rest-api/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var container *di.Container

func main() {
	log.Println("Initializing database...")
	db.InitDB()
	log.Println("Database initialized")

	// Initialize DI container
	log.Println("Initializing DI container...")
	container = di.NewContainer()
	log.Println("DI container initialized")

	// Start Kafka consumer in a goroutine
	log.Println("Starting Kafka consumer...")
	go startKafkaConsumer()

	// Start gRPC server in a goroutine
	log.Println("Starting gRPC server...")
	go startGRPCServer()

	// Start REST server
	log.Println("Starting REST server...")
	startRESTServer()
}

func startRESTServer() {
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

	// Initialize services in routes
	userService := container.GetUserService()
	eventService := container.GetEventService()
	authService := container.GetAuthService()
	routes.InitServices(userService, eventService, authService)

	routes.SetupRoutes(server)
	log.Println("REST server starting on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

func startGRPCServer() {
	log.Println("Creating gRPC listener...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC listener created successfully")

	grpcServer := grpc.NewServer()
	log.Println("gRPC server created")

	// Get services from DI container
	userService := container.GetUserService()
	eventService := container.GetEventService()
	authService := container.GetAuthService()

	// Create gRPC servers with DI
	authServer := auth.NewAuthServer(userService, authService)
	eventServer := event.NewEventServer(eventService)

	authpb.RegisterAuthServiceServer(grpcServer, authServer)
	eventpb.RegisterEventServiceServer(grpcServer, eventServer)
	log.Println("gRPC services registered")

	// Enable reflection for debugging
	reflection.Register(grpcServer)
	log.Println("gRPC reflection enabled")

	log.Println("gRPC server starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startKafkaConsumer() {
	consumer, err := kafka.NewConsumer()
	if err != nil {
		log.Printf("Failed to create Kafka consumer: %v", err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
	}()

	consumer.StartConsuming()
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
