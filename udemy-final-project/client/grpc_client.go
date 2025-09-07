// Package main provides a gRPC client for testing the event management API.
package main

import (
	"context"
	"log"
	"time"

	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/proto/auth"
	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	// Create clients
	authClient := auth.NewAuthServiceClient(conn)
	eventClient := event.NewEventServiceClient(conn)

	var token string

	// Register a new user
	log.Println("Registering user...")
	registerResp, err := authClient.Register(context.Background(), &auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Printf("Registration failed: %v", err)
	} else {
		log.Printf("User registered: %+v", registerResp.User)
	}

	// Login
	log.Println("Logging in...")
	loginResp, err := authClient.Login(context.Background(), &auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Printf("Login failed: %v", err)
	} else {
		log.Printf("Login successful: %s", loginResp.Message)
		log.Printf("Token: %s", loginResp.Token)
		token = loginResp.Token
	}

	// Create authenticated context for subsequent requests
	ctx := context.Background()
	if token != "" {
		md := metadata.New(map[string]string{"authorization": "Bearer " + token})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Create an event
	log.Println("Creating event...")
	eventResp, err := eventClient.CreateEvent(ctx, &event.CreateEventRequest{
		Name:        "gRPC Test Event",
		Description: "Testing gRPC event creation with Kafka",
		Location:    "Test Location",
		DateTime:    timestamppb.New(time.Now().Add(24 * time.Hour)),
	})
	if err != nil {
		log.Printf("Create event failed: %v", err)
	} else {
		log.Printf("Event created: %+v", eventResp.Event)
	}

	// Update the event
	log.Println("Updating event...")
	_, err = eventClient.UpdateEvent(ctx, &event.UpdateEventRequest{
		Id:          3, // Use the ID from the created event
		Name:        "Updated gRPC Test Event",
		Description: "Updated description for Kafka testing",
		Location:    "Updated Location",
		DateTime:    timestamppb.New(time.Now().Add(48 * time.Hour)),
	})
	if err != nil {
		log.Printf("Update event failed: %v", err)
	} else {
		log.Printf("Event updated successfully")
	}

	// Delete the event
	log.Println("Deleting event...")
	_, err = eventClient.DeleteEvent(ctx, &event.DeleteEventRequest{
		Id: 3,
	})
	if err != nil {
		log.Printf("Delete event failed: %v", err)
	} else {
		log.Printf("Event deleted successfully")
	}
}
