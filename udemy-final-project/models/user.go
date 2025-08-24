package models

import (
	"github.com/gurkanindibay/udemy-rest-api/db"
	"database/sql"
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

	result, err := stmt.Exec(u.Email, u.Password)
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
