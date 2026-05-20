package main

import (
	"fmt"
)

// ValidateSession checks if a sessionID corresponds to an active session.
func ValidateSession(sessionID string) (bool, string) {
	if sessionID == "" {
		return false, ""
	}

	username, err := DefaultAuth.ValidateSession(sessionID)
	if err != nil {
		return false, ""
	}

	return true, username
}

// GetUserInfo retrieves user information (requires active session)
func GetUserInfo(sessionID string) (*User, string) {
	user, err := DefaultAuth.GetUserInfo(sessionID)
	if err != nil {
		return nil, "Error: " + err.Error()
	}
	return user, ""
}

// ListAllUsers returns all registered users (for admin purposes)
func ListAllUsers() []string {
	return DefaultAuth.ListUsernames()
}

// DeleteUser removes a user from the system
func DeleteUser(username string) string {
	if err := DefaultAuth.DeleteUser(username); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ User '%s' deleted successfully\n", username)
	return ""
}

// ChangePassword allows a user to change their password (requires valid session)
func ChangePassword(sessionID, oldPassword, newPassword string) string {
	if err := DefaultAuth.ChangePassword(sessionID, oldPassword, newPassword); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ Password changed successfully for session '%s'\n", sessionID)
	return ""
}
