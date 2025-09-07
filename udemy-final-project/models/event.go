package models

import (
	"time"

	"github.com/gurkanindibay/udemy-rest-api/db"
	"gorm.io/gorm"
)

// Event represents an event in the system
type Event struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Name          string         `json:"name" gorm:"not null" binding:"required" example:"Sample Event"`
	Description   string         `json:"description" gorm:"not null" binding:"required" example:"This is a sample event"`
	Location      string         `json:"location" gorm:"not null" binding:"required" example:"Sample Location"`
	DateTime      time.Time      `json:"date_time" gorm:"not null" binding:"required" example:"2023-10-10T10:00:00Z"`
	UserID        int64          `json:"user_id,omitempty" gorm:"not null" example:"1"`
	User          User           `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Registrations []Registration `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"-"`
}

// Registration represents a user's registration for an event
type Registration struct {
	ID      int64 `gorm:"primaryKey;autoIncrement"`
	UserID  int64 `gorm:"not null"`
	EventID int64 `gorm:"not null"`
	User    User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Event   Event `gorm:"foreignKey:EventID;constraint:OnDelete:SET NULL" json:"-"`
}

// CreateEventRequest represents the request payload for creating an event
type CreateEventRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
}

// Save creates a new event in the database
func (e *Event) Save() error {
	gormDB := db.GetDB()
	return gormDB.Select("Name", "Description", "Location", "DateTime", "UserID").Create(e).Error
}

// GetAllEvents retrieves all events from the database
func GetAllEvents() ([]Event, error) {
	gormDB := db.GetDB()
	var events []Event
	err := gormDB.Find(&events).Error
	return events, err
}

// GetEventByID retrieves a specific event by its ID
func GetEventByID(id string) (*Event, error) {
	gormDB := db.GetDB()
	var event Event
	err := gormDB.First(&event, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

// Update modifies an existing event in the database
func (e *Event) Update() error {
	gormDB := db.GetDB()
	return gormDB.Model(e).Select("Name", "Description", "Location", "DateTime", "UserID").Updates(e).Error
}

// DeleteEvent removes an event from the database by its ID
func DeleteEvent(id string) error {
	gormDB := db.GetDB()
	return gormDB.Delete(&Event{}, id).Error
}

// Register creates a registration for a user to attend this event
func (e Event) Register(userID int64) error {
	gormDB := db.GetDB()
	registration := Registration{UserID: userID, EventID: e.ID}
	return gormDB.Create(&registration).Error
}

// CancelEventRegistration removes a user's registration for an event
func (e Event) CancelEventRegistration(userID int64, eventID string) error {
	gormDB := db.GetDB()
	return gormDB.Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&Registration{}).Error
}

// GetRegistrationsByUserID retrieves all events a user is registered for
func GetRegistrationsByUserID(userID int64) ([]Event, error) {
	gormDB := db.GetDB()
	var events []Event
	err := gormDB.Joins("JOIN registrations r ON events.id = r.event_id").
		Where("r.user_id = ?", userID).
		Find(&events).Error
	return events, err
}
