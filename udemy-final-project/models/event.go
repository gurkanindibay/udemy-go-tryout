package models

import (
	"time"

	"github.com/gurkanindibay/udemy-rest-api/db"
	"gorm.io/gorm"
)

type Event struct {
	ID          int64     `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" gorm:"not null" binding:"required" example:"Sample Event"`
	Description string    `json:"description" gorm:"not null" binding:"required" example:"This is a sample event"`
	Location    string    `json:"location" gorm:"not null" binding:"required" example:"Sample Location"`
	DateTime    time.Time `json:"date_time" gorm:"not null" binding:"required" example:"2023-10-10T10:00:00Z"`
	UserId      int64     `json:"user_id,omitempty" gorm:"not null" example:"1"`
	User        User      `gorm:"foreignKey:UserId"`
	Registrations []Registration `gorm:"foreignKey:EventId"`
}

type Registration struct {
	ID      int64 `gorm:"primaryKey"`
	UserId  int64 `gorm:"not null"`
	EventId int64 `gorm:"not null"`
	User    User  `gorm:"foreignKey:UserId"`
	Event   Event `gorm:"foreignKey:EventId"`
}

var events = []Event{}

func (e *Event) Save() error {
	gormDB := db.GetDB()
	return gormDB.Create(e).Error
}

func GetAllEvents() ([]Event, error) {
	gormDB := db.GetDB()
	var events []Event
	err := gormDB.Find(&events).Error
	return events, err
}

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

func (e *Event) Update() error {
	gormDB := db.GetDB()
	return gormDB.Save(e).Error
}

func DeleteEvent(id string) error {
	gormDB := db.GetDB()
	return gormDB.Delete(&Event{}, id).Error
}

func (e Event) Register(userId int64) error {
	gormDB := db.GetDB()
	registration := Registration{UserId: userId, EventId: e.ID}
	return gormDB.Create(&registration).Error
}

func (e Event) CancelEventRegistration(userId int64, eventId string) error {
	gormDB := db.GetDB()
	return gormDB.Where("user_id = ? AND event_id = ?", userId, eventId).Delete(&Registration{}).Error
}

func GetRegistrationsByUserID(userId int64) ([]Event, error) {
	gormDB := db.GetDB()
	var events []Event
	err := gormDB.Joins("JOIN registrations r ON events.id = r.event_id").
		Where("r.user_id = ?", userId).
		Find(&events).Error
	return events, err
}
