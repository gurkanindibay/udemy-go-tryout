package models

import (
	"time"

	"github.com/gurkanindibay/udemy-rest-api/db"
)

type Event struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" binding:"required" example:"Sample Event"`
	Description string    `json:"description" binding:"required" example:"This is a sample event"`
	Location    string    `json:"location" binding:"required" example:"Sample Location"`
	DateTime    time.Time `json:"date_time" binding:"required" example:"2023-10-10T10:00:00Z"`
	UserId      int64     `json:"user_id,omitempty" example:"1"`
}

type CreateEventRequest struct {
	Name        string    `json:"name" binding:"required" example:"Sample Event"`
	Description string    `json:"description" binding:"required" example:"This is a sample event"`
	Location    string    `json:"location" binding:"required" example:"Sample Location"`
	DateTime    time.Time `json:"date_time" binding:"required" example:"2023-10-10T10:00:00Z"`
}

var events = []Event{}

func (e Event) Save() error {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		INSERT INTO events (name, description, location, date_time, user_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(e.Name, e.Description, e.Location, e.DateTime, e.UserId).Scan(&e.ID)
	if err != nil {
		return err
	}

	events = append(events, e)
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `
		SELECT id, name, description, location, date_time, user_id
		FROM events
	`
	db := db.GetDB()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId); err != nil {
			panic("Failed to scan event: " + err.Error())
		}

		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(id string) (*Event, error) {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		SELECT id, name, description, location, date_time, user_id
		FROM events
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var e Event
	if err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId); err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *Event) Update() error {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		UPDATE events
		SET name = $1, description = $2, location = $3, date_time = $4, user_id = $5
		WHERE id = $6
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserId, e.ID)
	return err
}

func DeleteEvent(id string) error {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		DELETE FROM events
		WHERE id = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (e Event) Register(userId int64) error {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		INSERT INTO registrations (user_id, event_id)
		VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, e.ID)
	return err
}

func (e Event) CancelEventRegistration(userId int64, eventId string) error {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		DELETE FROM registrations
		WHERE user_id = $1 AND event_id = $2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, eventId)
	return err
}

func GetRegistrationsByUserID(userId int64) ([]Event, error) {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		SELECT e.id, e.name, e.description, e.location, e.date_time, e.user_id
		FROM events e
		JOIN registrations r ON e.id = r.event_id
		WHERE r.user_id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId); err != nil {
			panic("Failed to scan event: " + err.Error())
		}

		events = append(events, e)
	}

	return events, nil
}
