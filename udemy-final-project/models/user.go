package models

import (
	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/utils"
	"database/sql"
	"log"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (u *User) Save() error {
	db := db.GetDB()
	query := `
	INSERT INTO users (email, password) VALUES (?, ?);
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Log the email and password
	log.Printf("Attempting to register user: %s", u.Email)
	log.Printf("User password: %s", u.Password)


	// Hash the password before storing it
	hashedPassword, err := utils.HashPassword(u.Password)

	log.Printf("Hashed password: %s", hashedPassword)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	db := db.GetDB()
	query := `
	SELECT id, email, password FROM users WHERE email = ?;
	`
	row := db.QueryRow(query, email)

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
