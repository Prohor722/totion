package main

import (
	"fmt"
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

// LogoutUser ends a user session
func LogoutUser(sessionID string) string {
	session, exists := sessions[sessionID]
	if !exists {
		return "Error: Session not found"
	}

	session.IsActive = false
	delete(sessions, sessionID)
	fmt.Printf("✓ User '%s' logged out successfully\n", session.Username)
	return ""
}

