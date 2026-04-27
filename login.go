package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// User represents a user in the system
type User struct {
	Username     string
	Email        string
	PasswordHash string
}

// Session represents an active user session
type Session struct {
	Username  string
	SessionID string
	IsActive  bool
}

// UserDatabase stores all users (in-memory)
var userDatabase = make(map[string]*User)

// Sessions stores active sessions (in-memory)
var sessions = make(map[string]*Session)

// RegisterUser creates a new user account
// Returns error message if registration fails, empty string if successful


// LoginUser authenticates a user and creates a session
// Returns session ID if successful, error message if failed
func LoginUser(username, password string) (string, string) {
	if username == "" || password == "" {
		return "", "Error: Username and password are required"
	}

	// Check if user exists
	user, exists := userDatabase[username]
	if !exists {
		return "", "Error: Invalid username or password"
	}

	// Verify password
	if !verifyPassword(password, user.PasswordHash) {
		return "", "Error: Invalid username or password"
	}

	// Create session
	sessionID := generateSessionID(username)
	sessions[sessionID] = &Session{
		Username:  username,
		SessionID: sessionID,
		IsActive:  true,
	}

	fmt.Printf("✓ User '%s' logged in successfully\n", username)
	return sessionID, ""
}



// ValidateSession checks if a session is active
func ValidateSession(sessionID string) (bool, string) {
	session, exists := sessions[sessionID]
	if !exists || !session.IsActive {
		return false, ""
	}
	return true, session.Username
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
