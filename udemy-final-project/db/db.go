package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model for migration
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// Event model for migration
type Event struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Location    string    `gorm:"not null"`
	DateTime    time.Time `gorm:"not null"`
	UserId      int64     `gorm:"not null"`
}

// Registration model for migration
type Registration struct {
	ID      int64 `gorm:"primaryKey;autoIncrement"`
	UserId  int64 `gorm:"not null"`
	EventId int64 `gorm:"not null"`
}

var DB *gorm.DB

func InitDB() {
	var err error
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "eventdb")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	DB, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err = DB.AutoMigrate(&User{}, &Event{}, &Registration{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database connection established")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDB() *gorm.DB {
	return DB
}
