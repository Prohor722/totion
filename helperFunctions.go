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