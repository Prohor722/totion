package main

import (
	"crypto/sha256"
	"fmt"
)

// hashPassword creates a SHA256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// verifyPassword checks if the provided password matches the hash
func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

// generateSessionID creates a simple session ID
func generateSessionID(username string) string {
	hash := sha256.Sum256([]byte(username + fmt.Sprintf("%d", len(sessions))))
	return fmt.Sprintf("%x", hash)[:16]
}

// ValidateSession checks if a session is active
func ValidateSession(sessionID string) (bool, string) {
	session, exists := sessions[sessionID]
	if !exists || !session.IsActive {
		return false, ""
	}
	return true, session.Username
}


// GetUserInfo retrieves user information (requires active session)
func GetUserInfo(sessionID string) (*User, string) {
	isValid, username := ValidateSession(sessionID)
	if !isValid {
		return nil, "Error: Invalid or expired session"
	}

	user, exists := userDatabase[username]
	if !exists {
		return nil, "Error: User not found"
	}

	return user, ""
}

// ListAllUsers returns all registered users (for admin purposes)
func ListAllUsers() []string {
	var usernames []string
	for username := range userDatabase {
		usernames = append(usernames, username)
	}
	return usernames
}

// DeleteUser removes a user from the system
func DeleteUser(username string) string {
	if _, exists := userDatabase[username]; !exists {
		return "Error: User not found"
	}

	delete(userDatabase, username)
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

	user, exists := userDatabase[username]
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