package services

import (
	"errors"
	"strconv"

	"github.com/gurkanindibay/udemy-rest-api/kafka"
	"github.com/gurkanindibay/udemy-rest-api/models"
	"github.com/gurkanindibay/udemy-rest-api/utils"
)

// userServiceImpl implements UserService
type userServiceImpl struct{}

func NewUserService() UserService {
	return &userServiceImpl{}
}

func (s *userServiceImpl) Register(email, password string) (*models.User, error) {
	user := models.User{
		Email:    email,
		Password: password,
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userServiceImpl) Login(email, password string) (*models.User, error) {
	return models.VerifyUserCredentials(email, password)
}

// eventServiceImpl implements EventService
type eventServiceImpl struct{
	producer *kafka.Producer
}

func NewEventService() EventService {
	producer, err := kafka.NewProducer()
	if err != nil {
		// Log error but don't fail, allow service to work without Kafka
		producer = nil
	}
	return &eventServiceImpl{
		producer: producer,
	}
}

func (s *eventServiceImpl) GetAllEvents() ([]models.Event, error) {
	return models.GetAllEvents()
}

func (s *eventServiceImpl) GetEventByID(id string) (*models.Event, error) {
	return models.GetEventByID(id)
}

func (s *eventServiceImpl) CreateEvent(event models.Event) (*models.Event, error) {
	if err := event.Save(); err != nil {
		return nil, err
	}
	// Publish to Kafka
	if s.producer != nil {
		go s.producer.PublishEvent("created", strconv.FormatInt(event.ID, 10), event)
	}
	return &event, nil
}

func (s *eventServiceImpl) UpdateEvent(event models.Event) error {
	err := event.Update()
	if err == nil && s.producer != nil {
		go s.producer.PublishEvent("updated", strconv.FormatInt(event.ID, 10), event)
	}
	return err
}

func (s *eventServiceImpl) DeleteEvent(id string) error {
	event, err := s.GetEventByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}
	err = models.DeleteEvent(id)
	if err == nil && s.producer != nil {
		go s.producer.PublishEvent("deleted", strconv.FormatInt(event.ID, 10), *event)
	}
	return err
}

func (s *eventServiceImpl) RegisterForEvent(userID int64, eventID string) error {
	// Get the event to ensure it exists
	event, err := s.GetEventByID(eventID)
	if err != nil {
		return err
	}

	// Register the user for the event
	return event.Register(userID)
}

func (s *eventServiceImpl) CancelRegistration(userID int64, eventID string) error {
	event := models.Event{ID: 0} // We'll need to parse eventID
	return event.CancelEventRegistration(userID, eventID)
}

func (s *eventServiceImpl) GetUserRegistrations(userID int64) ([]models.Event, error) {
	return models.GetRegistrationsByUserID(userID)
}

// authServiceImpl implements AuthService
type authServiceImpl struct{}

func NewAuthService() AuthService {
	return &authServiceImpl{}
}

func (s *authServiceImpl) GenerateToken(email string, userID int64) (string, error) {
	return utils.GenerateToken(email, userID)
}

func (s *authServiceImpl) ValidateToken(tokenString string) (*models.User, error) {
	// This would need to be implemented in utils
	return nil, errors.New("not implemented")
}
