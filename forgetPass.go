package main

import "fmt"

// GetUserProfile retrieves the profile of a user using the profile service.
func GetUserProfile(username string) (*UserProfile, string) {
	profile, err := DefaultProfileService.GetProfile(username)
	if err != nil {
		return nil, "Error: " + err.Error()
	}
	return profile, ""
}

// UpdateUserProfile updates the profile information of a user.
func UpdateUserProfile(username, email, bio string) string {
	if err := DefaultProfileService.UpdateProfile(username, email, bio); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ User '%s' profile updated successfully\n", username)
	return ""
}
