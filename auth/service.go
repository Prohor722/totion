package auth

import (
	"sync"
	"time"

	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
)

// RegistrationService, CredentialService, SessionValidationService, PasswordService
// and AuthService compose the public auth API used by the application.
type RegistrationService interface {
	Register(username, email, password string) error
	DeleteUser(username string) error
	ListUsernames() []string
}

type CredentialService interface {
	Login(username, password string) (string, error)
	Logout(sessionID string) error
}

type SessionValidationService interface {
	ValidateSession(sessionID string) (string, error)
	GetUserInfo(sessionID string) (*model.User, error)
}

type PasswordService interface {
	ChangePassword(sessionID, oldPassword, newPassword string) error
	CreatePasswordResetToken(email string) (string, error)
	ResetPasswordWithToken(token, newPassword string) error
}

type AuthService interface {
	RegistrationService
	CredentialService
	SessionValidationService
	PasswordService
}

// authService is the concrete implementation behind the AuthService interface.
type authService struct {
	mutex         sync.RWMutex
	users         store.UserRepository
	sessions      SessionRepository
	resetTokens   map[string]passwordResetToken
	resetRequests map[string]passwordResetRequest
	failures      map[string]failedLogin
	clock         Clock
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// NewAuthService constructs an AuthService using provided repositories.
func NewAuthService(users store.UserRepository, sessions SessionRepository) AuthService {
	return NewAuthServiceWithClock(users, sessions, realClock{})
}

// NewAuthServiceWithClock constructs an AuthService with a testable clock.
func NewAuthServiceWithClock(users store.UserRepository, sessions SessionRepository, clock Clock) AuthService {
	return &authService{
		users:         users,
		sessions:      sessions,
		resetTokens:   make(map[string]passwordResetToken),
		resetRequests: make(map[string]passwordResetRequest),
		failures:      make(map[string]failedLogin),
		clock:         clock,
	}
}

// DefaultAuth is the package-level auth service used by the app.
var DefaultAuth AuthService = NewAuthService(store.UserStore, NewInMemorySessionRepository())
