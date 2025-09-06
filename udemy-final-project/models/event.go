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
	// save to database
	query := `
		INSERT INTO events (name, description, location, date_time, user_id)
		VALUES (?, ?, ?, ?, ?)
	`
	db := db.GetDB()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime.Format(time.RFC3339), e.UserId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = id
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
		var dateTimeStr string
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &dateTimeStr, &e.UserId); err != nil {
			panic("Failed to scan event: " + err.Error())
		}

		// Parse the datetime string back to time.Time
		parsedTime, err := time.Parse(time.RFC3339, dateTimeStr)
		if err != nil {
			panic("Failed to parse datetime: " + err.Error())
		}
		e.DateTime = parsedTime

		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(id string) (*Event, error) {
	query := `
		SELECT id, name, description, location, date_time, user_id
		FROM events
		WHERE id = ?
	`
	db := db.GetDB()

	row := db.QueryRow(query, id)

	var e Event
	var dateTimeStr string
	if err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &dateTimeStr, &e.UserId); err != nil {
		return nil, err
	}

	// Parse the datetime string back to time.Time
	parsedTime, err := time.Parse(time.RFC3339, dateTimeStr)
	if err != nil {
		return nil, err
	}
	e.DateTime = parsedTime

	return &e, nil
}

func (e *Event) Update() error {
	// update in database
	query := `
		UPDATE events
		SET name = ?, description = ?, location = ?, date_time = ?, user_id = ?
		WHERE id = ?
	`
	db := db.GetDB()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime.Format(time.RFC3339), e.UserId, e.ID)
	return err
}

func DeleteEvent(id string) error {
	// delete from database
	query := `
		DELETE FROM events
		WHERE id = ?
	`
	db := db.GetDB()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (e Event) Register(userId int64) error {
	// save to database
	query := `
		INSERT INTO registrations (user_id, event_id)
		VALUES (?, ?)
	`
	db := db.GetDB()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, e.ID)
	return err
}

func (e Event) CancelEventRegistration(userId int64, eventId string) error {
	// delete from database
	query := `
		DELETE FROM registrations
		WHERE user_id = ? AND event_id = ?
	`
	db := db.GetDB()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, eventId)
	return err
}

func GetRegistrationsByUserID(userId int64) ([]Event, error) {
	query := `
		SELECT e.id, e.name, e.description, e.location, e.date_time, e.user_id
		FROM events e
		JOIN registrations r ON e.id = r.event_id
		WHERE r.user_id = ?
	`
	db := db.GetDB()

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		var dateTimeStr string
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &dateTimeStr, &e.UserId); err != nil {
			panic("Failed to scan event: " + err.Error())
		}

		// Parse the datetime string back to time.Time
		parsedTime, err := time.Parse(time.RFC3339, dateTimeStr)
		if err != nil {
			panic("Failed to parse datetime: " + err.Error())
		}
		e.DateTime = parsedTime

		events = append(events, e)
	}

	return events, nil
}
