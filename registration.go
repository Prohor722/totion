package main

// AuthManager defines authentication operations.
// This interface follows Interface Segregation Principle - clients depend on specific operations.
type AuthManager interface {
	RegisterUser(username, email, password string) error
	LoginUser(username, password string) (string, error)
	LogoutUser(sessionID string) error
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
func (m *defaultAuthManager) RegisterUser(username, email, password string) error {
	return m.auth.Register(username, email, password)
}

// LoginUser authenticates a user and creates a session.
// This method is responsible for login only - Single Responsibility Principle.
func (m *defaultAuthManager) LoginUser(username, password string) (string, error) {
	return m.auth.Login(username, password)
}

// LogoutUser terminates a user session.
// This method is responsible for logout only - Single Responsibility Principle.
func (m *defaultAuthManager) LogoutUser(sessionID string) error {
	return m.auth.Logout(sessionID)
}

// Legacy package-level functions that use the default global AuthManager.
// These maintain backwards compatibility with existing code.

// authManager is the default global instance.
var authManager AuthManager = NewDefaultAuthManager(DefaultAuth)

// RegisterUser creates a new user account using the auth service.
func RegisterUser(username, email, password string) error {
	return authManager.RegisterUser(username, email, password)
}

// LoginUser authenticates a user and returns a session ID or error.
func LoginUser(username, password string) (string, error) {
	return authManager.LoginUser(username, password)
}

// LogoutUser terminates a user session.
func LogoutUser(sessionID string) error {
	return authManager.LogoutUser(sessionID)
}
