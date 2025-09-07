// Package event provides gRPC event service implementation.
package event

import (
	"context"
	"errors"
	"strconv"

	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/models"
	eventpb "github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/proto/event"
	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/security"
	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/services"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the gRPC EventService server
type Server struct {
	eventpb.UnimplementedEventServiceServer
	eventService services.EventService
}

// NewEventServer creates a new EventServer instance
func NewEventServer(eventService services.EventService) *Server {
	return &Server{
		eventService: eventService,
	}
}

// Helper function to extract user ID from gRPC context
func getUserIDFromContext(ctx context.Context) (int64, error) {
	// Get authorization header from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("no metadata found")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return 0, errors.New("authorization header is required")
	}

	tokenString := authHeader[0]
	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := security.ValidateToken(tokenString)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return userID, nil
}

// Helper function to convert model Event to protobuf Event
func convertToProtoEvent(e models.Event) *eventpb.Event {
	return &eventpb.Event{
		Id:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Location:    e.Location,
		DateTime:    timestamppb.New(e.DateTime),
		UserId:      e.UserID,
	}
}

// GetEvents retrieves all events via gRPC
func (s *Server) GetEvents(_ context.Context, _ *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
	events, err := s.eventService.GetAllEvents()
	if err != nil {
		return nil, err
	}

	var protoEvents []*eventpb.Event
	for _, event := range events {
		protoEvents = append(protoEvents, convertToProtoEvent(event))
	}

	return &eventpb.GetEventsResponse{
		Events: protoEvents,
	}, nil
}

// GetEvent retrieves a specific event by ID via gRPC
func (s *Server) GetEvent(_ context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	event, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}

	return &eventpb.GetEventResponse{
		Event: convertToProtoEvent(*event),
	}, nil
}

// CreateEvent creates a new event via gRPC
func (s *Server) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	event := models.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		DateTime:    req.DateTime.AsTime(),
		UserID:      userID,
	}

	createdEvent, err := s.eventService.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &eventpb.CreateEventResponse{
		Event: convertToProtoEvent(*createdEvent),
	}, nil
}

// UpdateEvent updates an existing event via gRPC
func (s *Server) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the event exists and belongs to the user
	existingEvent, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}
	if existingEvent.UserID != userID {
		return nil, errors.New("you do not have permission to update this event")
	}

	updatedEvent := models.Event{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		DateTime:    req.DateTime.AsTime(),
		UserID:      userID,
	}

	if err := s.eventService.UpdateEvent(updatedEvent); err != nil {
		return nil, err
	}

	return &eventpb.UpdateEventResponse{
		Event: convertToProtoEvent(updatedEvent),
	}, nil
}

// DeleteEvent deletes an event via gRPC
func (s *Server) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the event exists and belongs to the user
	event, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}
	if event.UserID != userID {
		return nil, errors.New("you do not have permission to delete this event")
	}

	if err := s.eventService.DeleteEvent(strconv.FormatInt(req.Id, 10)); err != nil {
		return nil, err
	}

	return &eventpb.DeleteEventResponse{}, nil
}

// RegisterForEvent registers a user for an event via gRPC
func (s *Server) RegisterForEvent(ctx context.Context, req *eventpb.RegisterForEventRequest) (*eventpb.RegisterForEventResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.eventService.RegisterForEvent(userID, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.RegisterForEventResponse{}, nil
}

// CancelRegistration cancels a user's registration for an event via gRPC
func (s *Server) CancelRegistration(ctx context.Context, req *eventpb.CancelRegistrationRequest) (*eventpb.CancelRegistrationResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.eventService.CancelRegistration(userID, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.CancelRegistrationResponse{}, nil
}

// GetUserRegistrations retrieves all events a user is registered for via gRPC
func (s *Server) GetUserRegistrations(ctx context.Context, _ *eventpb.GetUserRegistrationsRequest) (*eventpb.GetUserRegistrationsResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	events, err := s.eventService.GetUserRegistrations(userID)
	if err != nil {
		return nil, err
	}

	var protoEvents []*eventpb.Event
	for _, event := range events {
		protoEvents = append(protoEvents, convertToProtoEvent(event))
	}

	return &eventpb.GetUserRegistrationsResponse{
		Events: protoEvents,
	}, nil
}
