package models

import (
	"log"

	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/security"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Email    string `json:"email" gorm:"unique;not null" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" gorm:"not null" binding:"required,min=6" example:"password123"`
}

// Save creates a new user in the database with hashed password
func (u *User) Save() error {
	db := db.GetDB()

	// Log the email and password
	log.Printf("Attempting to register user: %s", u.Email)
	log.Printf("User password: %s", u.Password)

	// Hash the password before storing it
	hashedPassword, err := security.HashPassword(u.Password)
	if err != nil {
		return err
	}

	log.Printf("Hashed password: %s", hashedPassword)

	u.Password = hashedPassword

	return db.Create(u).Error
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*User, error) {
	db := db.GetDB()

	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// VerifyUserCredentials checks if the provided email and password match a user in the database
func VerifyUserCredentials(email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // User not found
	}

	// Compare the provided password with the stored hashed password
	if err := security.CheckPasswordHash(password, user.Password); err != nil {
		return nil, nil // Invalid password
	}

	return user, nil
}
