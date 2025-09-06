package event

import (
	"context"
	"errors"
	"strconv"

	"github.com/gurkanindibay/udemy-rest-api/models"
	eventpb "github.com/gurkanindibay/udemy-rest-api/proto/event"
	"github.com/gurkanindibay/udemy-rest-api/services"
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
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

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
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

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
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

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
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

	if err := s.eventService.RegisterForEvent(userId, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.RegisterForEventResponse{}, nil
}

func (s *EventServer) CancelRegistration(ctx context.Context, req *eventpb.CancelRegistrationRequest) (*eventpb.CancelRegistrationResponse, error) {
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

	if err := s.eventService.CancelRegistration(userId, strconv.FormatInt(req.EventId, 10)); err != nil {
		return nil, err
	}

	return &eventpb.CancelRegistrationResponse{}, nil
}

func (s *EventServer) GetUserRegistrations(ctx context.Context, req *eventpb.GetUserRegistrationsRequest) (*eventpb.GetUserRegistrationsResponse, error) {
	// Get user ID from context (would need authentication middleware)
	userId := int64(1) // Placeholder - should come from auth context

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
