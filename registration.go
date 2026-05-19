package main

import (
	"fmt"
)

// RegisterUser creates a new user account using the repository
// Returns error message if registration fails, empty string if successful
func RegisterUser(username, email, password string) string {
	if username == "" || email == "" || password == "" {
		return "Error: All fields are required"
	}

	if len(username) < 3 {
		return "Error: Username must be at least 3 characters"
	}

	if len(password) < 6 {
		return "Error: Password must be at least 6 characters"
	}

	if !isValidEmail(email) {
		return "Error: Invalid email format"
	}

	if UserStore.Exists(username) {
		return "Error: Username already exists"
	}

	if _, found := UserStore.FindByEmail(email); found {
		return "Error: Email already registered"
	}

	newUser := &User{
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password),
		Profile: &UserProfile{
			Username: username,
			Email:    email,
			Bio:      "",
		},
	}

	if err := UserStore.Add(newUser); err != nil {
		return "Error: " + err.Error()
	}

	fmt.Printf("✓ User '%s' registered successfully\n", username)
	return ""
}