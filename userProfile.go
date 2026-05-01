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

// GetUserProfile retrieves the profile of a user
func GetUserProfile(username string) (*UserProfile, string) {
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
func UpdateUserProfile(username, email, bio string) string {
	user, exists := userDatabase[username]
	if !exists {
		return "Error: User not found"
	}
	if email != "" {
		user.Email = email
	}
	if bio != "" {
		user.Bio = bio
	}
	// In a real application, you would also update the bio in the database
	fmt.Printf("✓ User '%s' profile updated successfully\n", username)
	return ""
}

