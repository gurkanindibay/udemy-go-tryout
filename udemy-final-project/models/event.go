package models

import (
	"time"

	"github.com/gurkanindibay/udemy-rest-api/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserId      int64     `json:"user_id"`
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
