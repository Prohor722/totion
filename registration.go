package main

import (
	"fmt"
)

// RegisterUser creates a new user account using the auth service.
// Returns error message if registration fails, empty string if successful.
func RegisterUser(username, email, password string) string {
	if err := DefaultAuth.Register(username, email, password); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ User '%s' registered successfully\n", username)
	return ""
}
