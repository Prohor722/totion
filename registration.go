package main

import (
	"fmt"
)

// AuthManager defines authentication operations.
// This interface follows Interface Segregation Principle - clients depend on specific operations.
type AuthManager interface {
	RegisterUser(username, email, password string) string
	LoginUser(username, password string) (string, string)
	LogoutUser(sessionID string) string
}

// defaultAuthManager implements AuthManager by delegating to an AuthService.
// This follows Dependency Inversion Principle - concrete implementation depends on abstraction.
type defaultAuthManager struct {
	auth AuthService
}

// NewDefaultAuthManager creates a new AuthManager with the given AuthService.
func NewDefaultAuthManager(auth AuthService) AuthManager {
	return &defaultAuthManager{auth: auth}
}

// RegisterUser creates a new user account.
// This method is responsible for user registration only - Single Responsibility Principle.
func (m *defaultAuthManager) RegisterUser(username, email, password string) string {
	if err := m.auth.Register(username, email, password); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ User '%s' registered successfully\n", username)
	return ""
}

// LoginUser authenticates a user and creates a session.
// This method is responsible for login only - Single Responsibility Principle.
func (m *defaultAuthManager) LoginUser(username, password string) (string, string) {
	sessionID, err := m.auth.Login(username, password)
	if err != nil {
		return "", "Error: " + err.Error()
	}
	return sessionID, ""
}

// LogoutUser terminates a user session.
// This method is responsible for logout only - Single Responsibility Principle.
func (m *defaultAuthManager) LogoutUser(sessionID string) string {
	if err := m.auth.Logout(sessionID); err != nil {
		return "Error: " + err.Error()
	}
	return ""
}

// Legacy package-level functions that use the default global AuthManager.
// These maintain backwards compatibility with existing code.

// authManager is the default global instance.
var authManager AuthManager = NewDefaultAuthManager(DefaultAuth)

// RegisterUser creates a new user account using the auth service.
// Returns error message if registration fails, empty string if successful.
func RegisterUser(username, email, password string) string {
	return authManager.RegisterUser(username, email, password)
}

// LoginUser authenticates a user and returns a session ID or error.
func LoginUser(username, password string) (string, string) {
	return authManager.LoginUser(username, password)
}

// LogoutUser terminates a user session.
func LogoutUser(sessionID string) string {
	return authManager.LogoutUser(sessionID)
}
