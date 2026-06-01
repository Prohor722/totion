package main

import "errors"

// UserManager defines operations for managing user accounts and sessions.
// This interface follows Interface Segregation Principle - clients depend on specific operations.
type UserManager interface {
	ValidateSession(sessionID string) (string, error)
	GetUserInfo(sessionID string) (*User, error)
	ListAllUsers() []string
	DeleteUser(username string) error
	ChangePassword(sessionID, oldPassword, newPassword string) error
}

// defaultUserManager implements UserManager by delegating to focused auth abstractions.
// This follows Dependency Inversion Principle - concrete implementation depends on small interfaces.
type defaultUserManager struct {
	session SessionValidationService
	account RegistrationService
	password PasswordService
}

// NewDefaultUserManager creates a new UserManager with the given AuthService.
func NewDefaultUserManager(auth AuthService) UserManager {
	return &defaultUserManager{
		session: auth,
		account: auth,
		password: auth,
	}
}

// ValidateSession checks if a sessionID corresponds to an active session.
func (m *defaultUserManager) ValidateSession(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("session ID is required")
	}

	return m.session.ValidateSession(sessionID)
}

// GetUserInfo retrieves user information (requires active session).
func (m *defaultUserManager) GetUserInfo(sessionID string) (*User, error) {
	return m.session.GetUserInfo(sessionID)
}

// ListAllUsers returns all registered users (for admin purposes).
func (m *defaultUserManager) ListAllUsers() []string {
	return m.account.ListUsernames()
}

// DeleteUser removes a user from the system.
// This method is responsible for user deletion only - Single Responsibility Principle.
func (m *defaultUserManager) DeleteUser(username string) error {
	return m.account.DeleteUser(username)
}

// ChangePassword allows a user to change their password (requires valid session).
// This method is responsible for password changes only - Single Responsibility Principle.
func (m *defaultUserManager) ChangePassword(sessionID, oldPassword, newPassword string) error {
	return m.password.ChangePassword(sessionID, oldPassword, newPassword)
}

// Legacy package-level functions that use the default global UserManager.
// These maintain backwards compatibility with existing code.

// userManager is the default global instance.
var userManager UserManager = NewDefaultUserManager(DefaultAuth)

// ValidateSession checks if a sessionID corresponds to an active session.
func ValidateSession(sessionID string) (string, error) {
	return userManager.ValidateSession(sessionID)
}

// GetUserInfo retrieves user information (requires active session)
func GetUserInfo(sessionID string) (*User, error) {
	return userManager.GetUserInfo(sessionID)
}

// ListAllUsers returns all registered users (for admin purposes)
func ListAllUsers() []string {
	return userManager.ListAllUsers()
}

// DeleteUser removes a user from the system
func DeleteUser(username string) error {
	return userManager.DeleteUser(username)
}

// ChangePassword allows a user to change their password (requires valid session)
func ChangePassword(sessionID, oldPassword, newPassword string) error {
	return userManager.ChangePassword(sessionID, oldPassword, newPassword)
}
