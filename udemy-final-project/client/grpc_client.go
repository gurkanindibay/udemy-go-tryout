package main

import (
	"context"
	"log"
	"time"

	"github.com/gurkanindibay/udemy-rest-api/proto/auth"
	"github.com/gurkanindibay/udemy-rest-api/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	authClient := auth.NewAuthServiceClient(conn)
	eventClient := event.NewEventServiceClient(conn)

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
	}

	// Create an event
	log.Println("Creating event...")
	eventResp, err := eventClient.CreateEvent(context.Background(), &event.CreateEventRequest{
		Name:        "gRPC Test Event",
		Description: "Testing gRPC event creation",
		Location:    "Test Location",
		DateTime:    timestamppb.New(time.Now().Add(24 * time.Hour)),
	})
	if err != nil {
		log.Printf("Create event failed: %v", err)
	} else {
		log.Printf("Event created: %+v", eventResp.Event)
	}

	// Get all events
	log.Println("Getting all events...")
	eventsResp, err := eventClient.GetEvents(context.Background(), &event.GetEventsRequest{})
	if err != nil {
		log.Printf("Get events failed: %v", err)
	} else {
		log.Printf("Found %d events", len(eventsResp.Events))
		for _, e := range eventsResp.Events {
			log.Printf("Event: %s - %s", e.Name, e.Description)
		}
	}
}
