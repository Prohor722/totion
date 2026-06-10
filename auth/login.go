package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
	"github.com/Prohor722/totion/util"
)

type SessionRepository interface {
	Create(session *model.Session)
	Get(sessionID string) (*model.Session, bool)
	Remove(sessionID string)
}

type InMemorySessionRepository struct {
	sessions map[string]*model.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]*model.Session)}
}

func (r *InMemorySessionRepository) Create(session *model.Session) {
	if session == nil {
		return
	}
	r.sessions[session.SessionID] = session
}

func (r *InMemorySessionRepository) Get(sessionID string) (*model.Session, bool) {
	session, ok := r.sessions[sessionID]
	return session, ok
}

func (r *InMemorySessionRepository) Remove(sessionID string) {
	delete(r.sessions, sessionID)
}

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

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

type authService struct {
	users       store.UserRepository
	sessions    SessionRepository
	resetTokens map[string]passwordResetToken
	failures    map[string]failedLogin
	clock       Clock
}

type passwordResetToken struct {
	Username  string
	ExpiresAt time.Time
}

type failedLogin struct {
	Count       int
	LastAttempt time.Time
	LockedUntil time.Time
}

func NewAuthService(users store.UserRepository, sessions SessionRepository) AuthService {
	return NewAuthServiceWithClock(users, sessions, realClock{})
}

func NewAuthServiceWithClock(users store.UserRepository, sessions SessionRepository, clock Clock) AuthService {
	return &authService{
		users:       users,
		sessions:    sessions,
		resetTokens: make(map[string]passwordResetToken),
		failures:    make(map[string]failedLogin),
		clock:       clock,
	}
}

// DefaultAuth is the package-level auth service used by the app.
var DefaultAuth AuthService = NewAuthService(store.UserStore, NewInMemorySessionRepository())

func (s *authService) Register(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("all fields are required")
	}

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if !util.IsStrongPassword(password) {
		return fmt.Errorf("password must be at least %d characters and include upper, lower, digit, and symbol", util.MinPasswordLength)
	}

	if !util.IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	if s.users.Exists(username) {
		return errors.New("username already exists")
	}

	if _, found := s.users.FindByEmail(email); found {
		return errors.New("email already registered")
	}

	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	return s.users.Add(&model.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Profile: &model.UserProfile{
			Username: username,
			Email:    email,
			Bio:      "",
		},
	})
}

func (s *authService) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	if s.isAccountLocked(username) {
		return "", errors.New("account temporarily locked due to failed login attempts")
	}

	user, exists := s.users.Get(username)
	if !exists || !util.VerifyPassword(password, user.PasswordHash) {
		s.recordFailedLogin(username)
		return "", errors.New("invalid username or password")
	}

	s.clearFailedLogin(username)

	sessionID := util.GenerateSessionID(username)
	s.sessions.Create(&model.Session{
		Username:  username,
		SessionID: sessionID,
		IsActive:  true,
		ExpiresAt: s.clock.Now().Add(util.SessionTTL),
	})

	return sessionID, nil
}

func (s *authService) Logout(sessionID string) error {
	session, exists := s.sessions.Get(sessionID)
	if !exists || !session.IsActive {
		return errors.New("session not found")
	}

	session.IsActive = false
	s.sessions.Remove(sessionID)
	return nil
}

func (s *authService) ValidateSession(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("session ID is required")
	}

	session, exists := s.sessions.Get(sessionID)
	if !exists || !session.IsActive {
		return "", errors.New("invalid or expired session")
	}

	if s.clock.Now().After(session.ExpiresAt) {
		s.sessions.Remove(sessionID)
		return "", errors.New("session expired")
	}

	return session.Username, nil
}

func (s *authService) isAccountLocked(username string) bool {
	failure, exists := s.failures[username]
	if !exists {
		return false
	}
	return s.clock.Now().Before(failure.LockedUntil)
}

func (s *authService) recordFailedLogin(username string) {
	failure := s.failures[username]
	now := s.clock.Now()
	if now.Sub(failure.LastAttempt) > util.FailedLoginWindow {
		failure.Count = 0
	}
	failure.Count++
	failure.LastAttempt = now
	if failure.Count >= util.FailedLoginThreshold {
		failure.LockedUntil = now.Add(util.AccountLockDuration)
	}
	s.failures[username] = failure
}

func (s *authService) clearFailedLogin(username string) {
	delete(s.failures, username)
}

func (s *authService) GetUserInfo(sessionID string) (*model.User, error) {
	username, err := s.ValidateSession(sessionID)
	if err != nil {
		return nil, err
	}

	user, exists := s.users.Get(username)
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) ChangePassword(sessionID, oldPassword, newPassword string) error {
	username, err := s.ValidateSession(sessionID)
	if err != nil {
		return err
	}

	user, exists := s.users.Get(username)
	if !exists {
		return errors.New("user not found")
	}

	if !util.VerifyPassword(oldPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if !util.IsStrongPassword(newPassword) {
		return fmt.Errorf("new password must be at least %d characters and include upper, lower, digit, and symbol", util.MinPasswordLength)
	}

	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	return nil
}

func (s *authService) CreatePasswordResetToken(email string) (string, error) {
	user, exists := s.users.FindByEmail(email)
	if !exists {
		return "", errors.New("email not found")
	}
	token := util.GenerateSessionID(user.Username)
	s.resetTokens[token] = passwordResetToken{Username: user.Username, ExpiresAt: s.clock.Now().Add(util.PasswordResetTokenTTL)}
	return token, nil
}

func (s *authService) ResetPasswordWithToken(token, newPassword string) error {
	reset, exists := s.resetTokens[token]
	if !exists {
		return errors.New("invalid token")
	}

	if s.clock.Now().After(reset.ExpiresAt) {
		delete(s.resetTokens, token)
		return errors.New("reset token expired")
	}

	if !util.IsStrongPassword(newPassword) {
		return fmt.Errorf("new password must be at least %d characters and include upper, lower, digit, and symbol", util.MinPasswordLength)
	}

	user, exists := s.users.Get(reset.Username)
	if !exists {
		return errors.New("user not found")
	}

	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	delete(s.resetTokens, token)
	return nil
}

func (s *authService) ListUsernames() []string {
	return s.users.List()
}

func (s *authService) DeleteUser(username string) error {
	return s.users.Delete(username)
}
