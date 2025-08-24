package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		panic("Failed to open database: " + err.Error())
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	createTables()
	fmt.Println("Database connection established")
}

func GetDB() *sql.DB {
	return db
}

func createTables() {
	createEventsTable()
}

func createEventsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		date_time TEXT NOT NULL,
		user_id INTEGER NOT NULL
	);
	`
	if _, err := db.Exec(query); err != nil {
		panic("Failed to create events table: " + err.Error())
	}
}
