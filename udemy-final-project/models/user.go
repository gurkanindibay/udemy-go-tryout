package models

import (
	"database/sql"
	"log"

	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/utils"
)

type User struct {
	ID       int64  `json:"id" example:"1"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

func (u *User) Save() error {
	db := db.GetDB()

	// Log the email and password
	log.Printf("Attempting to register user: %s", u.Email)
	log.Printf("User password: %s", u.Password)

	// Hash the password before storing it
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	log.Printf("Hashed password: %s", hashedPassword)

	// Use prepared statement for extra security
	stmt, err := db.Prepare(`
		INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Email, hashedPassword).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	db := db.GetDB()

	// Use prepared statement
	stmt, err := db.Prepare(`
		SELECT id, email, password FROM users WHERE email = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &u, nil
}

func VerifyUserCredentials(email, password string) (*User, error) {

	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // User not found
	}

	// Compare the provided password with the stored hashed password
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return nil, nil // Invalid password
	}

	return user, nil
}
