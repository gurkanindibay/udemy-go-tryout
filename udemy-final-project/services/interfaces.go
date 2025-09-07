package services

import (
	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/models"
)

// UserService interface for user operations
type UserService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
}

// EventService interface for event operations
type EventService interface {
	GetAllEvents() ([]models.Event, error)
	GetEventByID(id string) (*models.Event, error)
	CreateEvent(event models.Event) (*models.Event, error)
	UpdateEvent(event models.Event) error
	DeleteEvent(id string) error
	RegisterForEvent(userID int64, eventID string) error
	CancelRegistration(userID int64, eventID string) error
	GetUserRegistrations(userID int64) ([]models.Event, error)
}

// AuthService interface for authentication operations
type AuthService interface {
	GenerateToken(email string, userID int64) (string, error)
	ValidateToken(tokenString string) (*models.User, error)
}
