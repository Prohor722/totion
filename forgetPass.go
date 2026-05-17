package main

import (
	"fmt"
)

// UserProfile represents a user's profile information
type UserProfile struct {
	Username string
	Email    string
	Bio      string
}

// User represents a stored user record
type User struct {
	Username string
	Email    string
}

// userDatabase is an in-memory mock of a user store
var userDatabase = map[string]*User{}

// GetUserProfile retrieves the profile of a user
func GetUserProfile2(username string) (*UserProfile, string) {
	user, exists := userDatabase[username]
	if !exists {
		return nil, "Error: User not found"
	}
	profile := &UserProfile{
		Username: user.Username,
		Email:    user.Email,
		Bio:      "This is a user bio.", // Placeholder bio
	}
	return profile, ""
}

// UpdateUserProfile updates the profile information of a user
func UpdateUserProfile2(username, email, bio string) string {
	user, exists := userDatabase[username]
	if !exists {
		return "Error: User not found"
	}
	if email != "" {
		user.Email = email
	}
	// In a real application, you would also update the bio in the database
	fmt.Printf("✓ User '%s' profile updated successfully\n", username)
	return ""
}