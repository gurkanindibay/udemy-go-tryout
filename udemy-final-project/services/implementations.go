package services

import (
	"errors"

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
type eventServiceImpl struct{}

func NewEventService() EventService {
	return &eventServiceImpl{}
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
	return &event, nil
}

func (s *eventServiceImpl) UpdateEvent(event models.Event) error {
	return event.Update()
}

func (s *eventServiceImpl) DeleteEvent(id string) error {
	return models.DeleteEvent(id)
}

func (s *eventServiceImpl) RegisterForEvent(userID int64, eventID string) error {
	event := models.Event{ID: 0} // We'll need to parse eventID
	// This is a simplified version - in real implementation you'd need to get the event first
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
