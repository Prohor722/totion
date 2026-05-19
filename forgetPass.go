package main

import "fmt"

// GetUserProfile2 retrieves the profile of a user using the repository
func GetUserProfile2(username string) (*UserProfile, string) {
	user, exists := UserStore.Get(username)
	if !exists {
		return nil, "Error: User not found"
	}
	// copy profile data
	profile := &UserProfile{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Profile.Bio,
	}
	return profile, ""
}

// UpdateUserProfile2 updates the profile information of a user
func UpdateUserProfile2(username, email, bio string) string {
	user, exists := UserStore.Get(username)
	if !exists {
		return "Error: User not found"
	}
	if email != "" {
		user.Email = email
		user.Profile.Email = email
	}
	if bio != "" {
		user.Profile.Bio = bio
	}
	fmt.Printf("✓ User '%s' profile updated successfully\n", username)
	return ""
}