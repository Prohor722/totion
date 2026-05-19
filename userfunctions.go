package main

import (
	"fmt"
)

// ValidateSession checks if a sessionID corresponds to an active session.
// For now, treat a non-empty sessionID as valid and return it as the username.
func ValidateSession(sessionID string) (bool, string) {
	if sessionID == "" {
		return false, ""
	}
	// In a real implementation, sessionID would be looked up in a session store.
	return true, sessionID
}

// GetUserInfo retrieves user information (requires active session)
func GetUserInfo(sessionID string) (*User, string) {
	isValid, username := ValidateSession(sessionID)
	if !isValid {
		return nil, "Error: Invalid or expired session"
	}

	user, exists := UserStore.Get(username)
	if !exists {
		return nil, "Error: User not found"
	}

	return user, ""
}

// ListAllUsers returns all registered users (for admin purposes)
func ListAllUsers() []string {
	return UserStore.List()
}

// DeleteUser removes a user from the system
func DeleteUser(username string) string {
	if err := UserStore.Delete(username); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ User '%s' deleted successfully\n", username)
	return ""
}

// ChangePassword allows a user to change their password (requires valid session)
func ChangePassword(sessionID, oldPassword, newPassword string) string {
	isValid, username := ValidateSession(sessionID)
	if !isValid {
		return "Error: Invalid or expired session"
	}

	if newPassword == "" || len(newPassword) < 6 {
		return "Error: New password must be at least 6 characters"
	}

	user, exists := UserStore.Get(username)
	if !exists {
		return "Error: User not found"
	}

	if !verifyPassword(oldPassword, user.PasswordHash) {
		return "Error: Current password is incorrect"
	}

	user.PasswordHash = hashPassword(newPassword)
	fmt.Printf("✓ Password changed successfully for user '%s'\n", username)
	return ""
}