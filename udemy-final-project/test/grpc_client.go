// Package main provides a comprehensive gRPC client test suite for the event management API.
package main

import (
	"context"
	"fmt"
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

	// Test variables
	testEmail := "grpc-test@example.com"
	testPassword := "grpc-test-password" // #nosec G101

	fmt.Println("Starting gRPC tests...")

	// Test 1: Register user
	fmt.Println("Test 1: gRPC User Registration")
	registerReq := &auth.RegisterRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	registerResp, err := authClient.Register(context.Background(), registerReq)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	fmt.Printf("Registration successful: User ID %d\n", registerResp.User.Id)

	// Test 2: Login user
	fmt.Println("Test 2: gRPC User Login")
	loginReq := &auth.LoginRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	loginResp, err := authClient.Login(context.Background(), loginReq)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Printf("Login successful: Token received\n")
	jwtToken := loginResp.Token

	// Helper function to create authenticated context
	createAuthContext := func() context.Context {
		return metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+jwtToken)
	}

	// Test 3: Get Events (should be empty)
	fmt.Println("Test 3: gRPC Get Events (empty)")
	getEventsReq := &event.GetEventsRequest{}
	getEventsResp, err := eventClient.GetEvents(context.Background(), getEventsReq)
	if err != nil {
		fmt.Printf("Get events failed: %v\n", err)
		return
	}
	if len(getEventsResp.Events) == 0 {
		fmt.Println("Get events successful: Empty list as expected")
	} else {
		fmt.Printf("Get events failed: Expected empty list, got %d events\n", len(getEventsResp.Events))
		return
	}

	// Test 4: Create Event
	fmt.Println("Test 4: gRPC Create Event")
	createEventReq := &event.CreateEventRequest{
		Name:        "gRPC Test Event",
		Description: "This is a test event created via gRPC",
		Location:    "gRPC Test Location",
		DateTime:    timestamppb.New(time.Now().Add(24 * time.Hour)),
	}

	createEventResp, err := eventClient.CreateEvent(createAuthContext(), createEventReq)
	if err != nil {
		fmt.Printf("Create event failed: %v\n", err)
		return
	}
	eventID := createEventResp.Event.Id
	fmt.Printf("Create event successful: Event ID %d\n", eventID)

	// Test 5: Get Events (should have 1 event)
	fmt.Println("Test 5: gRPC Get Events (should return 1 event)")
	getEventsResp, err = eventClient.GetEvents(context.Background(), getEventsReq)
	if err != nil {
		fmt.Printf("Get events failed: %v\n", err)
		return
	}
	if len(getEventsResp.Events) == 1 {
		fmt.Println("Get events successful: 1 event returned")
	} else {
		fmt.Printf("Get events failed: Expected 1 event, got %d events\n", len(getEventsResp.Events))
		return
	}

	// Test 6: Get Event by ID
	fmt.Println("Test 6: gRPC Get Event by ID")
	getEventReq := &event.GetEventRequest{
		Id: eventID,
	}

	getEventResp, err := eventClient.GetEvent(context.Background(), getEventReq)
	if err != nil {
		fmt.Printf("Get event by ID failed: %v\n", err)
		return
	}
	if getEventResp.Event.Id == eventID {
		fmt.Println("Get event by ID successful")
	} else {
		fmt.Printf("Get event by ID failed: Expected ID %d, got %d\n", eventID, getEventResp.Event.Id)
		return
	}

	// Test 7: Update Event
	fmt.Println("Test 7: gRPC Update Event")
	updateEventReq := &event.UpdateEventRequest{
		Id:          eventID,
		Name:        "Updated gRPC Test Event",
		Description: "This is an updated test event via gRPC",
		Location:    "Updated gRPC Test Location",
		DateTime:    timestamppb.New(time.Now().Add(48 * time.Hour)),
	}

	updateEventResp, err := eventClient.UpdateEvent(createAuthContext(), updateEventReq)
	if err != nil {
		fmt.Printf("Update event failed: %v\n", err)
		return
	}
	if updateEventResp.Event.Name == "Updated gRPC Test Event" {
		fmt.Println("Update event successful")
	} else {
		fmt.Println("Update event failed: Name not updated correctly")
		return
	}

	// Test 8: Register for Event
	fmt.Println("Test 8: gRPC Register for Event")
	registerForEventReq := &event.RegisterForEventRequest{
		EventId: eventID,
	}

	_, err = eventClient.RegisterForEvent(createAuthContext(), registerForEventReq)
	if err != nil {
		fmt.Printf("Register for event failed: %v\n", err)
		return
	}
	fmt.Println("Register for event successful")

	// Add a small delay to ensure registration is committed
	time.Sleep(1 * time.Second)

	// Test 9: Get User Registrations
	fmt.Println("Test 9: gRPC Get User Registrations")
	getUserRegistrationsReq := &event.GetUserRegistrationsRequest{}

	getUserRegistrationsResp, err := eventClient.GetUserRegistrations(createAuthContext(), getUserRegistrationsReq)
	if err != nil {
		fmt.Printf("Get user registrations failed: %v\n", err)
		return
	}
	if len(getUserRegistrationsResp.Events) >= 1 {
		fmt.Println("Get user registrations successful")
	} else {
		fmt.Println("Get user registrations failed: Expected at least 1 registration")
		return
	}

	// Test 10: Cancel Registration
	fmt.Println("Test 10: gRPC Cancel Registration")
	cancelRegistrationReq := &event.CancelRegistrationRequest{
		EventId: eventID,
	}

	_, err = eventClient.CancelRegistration(createAuthContext(), cancelRegistrationReq)
	if err != nil {
		fmt.Printf("Cancel registration failed: %v\n", err)
		return
	}
	fmt.Println("Cancel registration successful")

	// Test 11: Delete Event
	fmt.Println("Test 11: gRPC Delete Event")
	deleteEventReq := &event.DeleteEventRequest{
		Id: eventID,
	}

	_, err = eventClient.DeleteEvent(createAuthContext(), deleteEventReq)
	if err != nil {
		fmt.Printf("Delete event failed: %v\n", err)
		return
	}
	fmt.Println("Delete event successful")

	fmt.Println("All gRPC tests passed! ✅")
}
