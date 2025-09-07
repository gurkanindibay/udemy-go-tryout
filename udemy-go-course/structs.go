package main

import (
	"fmt"

	"github.com/gurkanindibay/udemy-go-tryout/udemy-go-course/user"
)

func test_structs() {
	user, err := user.NewUser(
		getUserData("Please enter your first name: "),
		getUserData("Please enter your last name: "),
		getUserData("Please enter your birthdate (MM/DD/YYYY): "),
	)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	outputUserData(*user)
	user.ClearUserData()

	fmt.Println("User data cleared.")
	fmt.Println("After clearing:")
	fmt.Println(user.OutputUserDetails())
}

func outputUserData(user user.User) {
	fmt.Println("User Information:")
	fmt.Println("First Name:", user.FirstName)
	fmt.Println(user.OutputUserDetails())
}

func getUserData(promptText string) string {
	fmt.Print(promptText)
	var value string
	_, err := fmt.Scanln(&value)
	if err != nil {
		return ""
	}
	return value
}
