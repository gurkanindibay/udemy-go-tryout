package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "eventdb")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", dataSourceName)
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDB() *sql.DB {
	return db
}

func createTables() {
	createUsersTable()
	createEventsTable()
	createRegistrationsTable()
}

func createEventsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		date_time TIMESTAMP NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id)
	);
	`
	if _, err := db.Exec(query); err != nil {
		panic("Failed to create events table: " + err.Error())
	}
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	if _, err := db.Exec(query); err != nil {
		panic("Failed to create users table: " + err.Error())
	}
}

func createRegistrationsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS registrations (
		id SERIAL PRIMARY KEY,
		event_id INTEGER NOT NULL REFERENCES events(id),
		user_id INTEGER NOT NULL REFERENCES users(id),
		UNIQUE(event_id, user_id)
	);
	`
	if _, err := db.Exec(query); err != nil {
		panic("Failed to create registrations table: " + err.Error())
	}
}
