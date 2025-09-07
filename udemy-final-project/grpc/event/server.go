package event

import (
	"context"
	"errors"
	"strconv"

	"github.com/gurkanindibay/udemy-rest-api/models"
	eventpb "github.com/gurkanindibay/udemy-rest-api/proto/event"
	"github.com/gurkanindibay/udemy-rest-api/services"
	"github.com/gurkanindibay/udemy-rest-api/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServer struct {
	eventpb.UnimplementedEventServiceServer
	eventService services.EventService
}

func NewEventServer(eventService services.EventService) *EventServer {
	return &EventServer{
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

	userId, err := utils.ValidateToken(tokenString)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return userId, nil
}

// Helper function to convert model Event to protobuf Event
func convertToProtoEvent(e models.Event) *eventpb.Event {
	return &eventpb.Event{
		Id:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Location:    e.Location,
		DateTime:    timestamppb.New(e.DateTime),
		UserId:      e.UserId,
	}
}

// Helper function to convert protobuf Event to model Event
func convertFromProtoEvent(pe *eventpb.Event) models.Event {
	return models.Event{
		ID:          pe.Id,
		Name:        pe.Name,
		Description: pe.Description,
		Location:    pe.Location,
		DateTime:    pe.DateTime.AsTime(),
		UserId:      pe.UserId,
	}
}

func (s *EventServer) GetEvents(ctx context.Context, req *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
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

func (s *EventServer) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	event, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}

	return &eventpb.GetEventResponse{
		Event: convertToProtoEvent(*event),
	}, nil
}

func (s *EventServer) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	event := models.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		DateTime:    req.DateTime.AsTime(),
		UserId:      userId,
	}

	createdEvent, err := s.eventService.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &eventpb.CreateEventResponse{
		Event: convertToProtoEvent(*createdEvent),
	}, nil
}

func (s *EventServer) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the event exists and belongs to the user
	existingEvent, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}
	if existingEvent.UserId != userId {
		return nil, errors.New("you do not have permission to update this event")
	}

	updatedEvent := models.Event{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		DateTime:    req.DateTime.AsTime(),
		UserId:      userId,
	}

	if err := s.eventService.UpdateEvent(updatedEvent); err != nil {
		return nil, err
	}

	return &eventpb.UpdateEventResponse{
		Event: convertToProtoEvent(updatedEvent),
	}, nil
}

func (s *EventServer) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the event exists and belongs to the user
	event, err := s.eventService.GetEventByID(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return nil, err
	}
	if event.UserId != userId {
		return nil, errors.New("you do not have permission to delete this event")
	}

	if err := s.eventService.DeleteEvent(strconv.FormatInt(req.Id, 10)); err != nil {
		return nil, err
	}

	return &eventpb.DeleteEventResponse{}, nil
}

func (s *EventServer) RegisterForEvent(ctx context.Context, req *eventpb.RegisterForEventRequest) (*eventpb.RegisterForEventResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.eventService.RegisterForEvent(userId, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.RegisterForEventResponse{}, nil
}

func (s *EventServer) CancelRegistration(ctx context.Context, req *eventpb.CancelRegistrationRequest) (*eventpb.CancelRegistrationResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.eventService.CancelRegistration(userId, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.CancelRegistrationResponse{}, nil
}

func (s *EventServer) GetUserRegistrations(ctx context.Context, req *eventpb.GetUserRegistrationsRequest) (*eventpb.GetUserRegistrationsResponse, error) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	events, err := s.eventService.GetUserRegistrations(userId)
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
