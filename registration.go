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
