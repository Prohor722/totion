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
