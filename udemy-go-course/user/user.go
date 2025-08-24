package user

import (
	"errors"
	"fmt"
	"time"
)

type User struct{
	FirstName string
	lastName  string
	birthdate string
	createdAt time.Time
}

func NewUser(firstName, lastName, birthdate string) (*User, error) {

	if firstName == "" || lastName == "" || birthdate == "" {
		return nil, errors.New("all fields must be filled")
	}

	return &User{
		FirstName: firstName,
		lastName:  lastName,
		birthdate: birthdate,
		createdAt: time.Now(),
	}, nil
}

func (user User) OutputUserDetails() string {
	return fmt.Sprintf("First Name: %s\nLast Name: %s\nBirthdate: %s\nCreated At: %s",
		user.FirstName, user.lastName, user.birthdate, user.createdAt)
}

func (user *User) ClearUserData() {
	user.FirstName = ""
	user.lastName = ""
}