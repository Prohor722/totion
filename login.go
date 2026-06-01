package main

import (
	"errors"
	"time"
)

type SessionRepository interface {
	Create(session *Session)
	Get(sessionID string) (*Session, bool)
	Remove(sessionID string)
}

type InMemorySessionRepository struct {
	sessions map[string]*Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]*Session)}
}

func (r *InMemorySessionRepository) Create(session *Session) {
	if session == nil {
		return
	}
	r.sessions[session.SessionID] = session
}

func (r *InMemorySessionRepository) Get(sessionID string) (*Session, bool) {
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
	GetUserInfo(sessionID string) (*User, error)
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

type authService struct {
	users       UserRepository
	sessions    SessionRepository
	resetTokens map[string]passwordResetToken
	failures    map[string]failedLogin
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

func NewAuthService(users UserRepository, sessions SessionRepository) AuthService {
	return &authService{
		users:       users,
		sessions:    sessions,
		resetTokens: make(map[string]passwordResetToken),
		failures:    make(map[string]failedLogin),
	}
}

// DefaultAuth is the package-level auth service used by the app.
var DefaultAuth AuthService = NewAuthService(UserStore, NewInMemorySessionRepository())

func (s *authService) Register(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("all fields are required")
	}

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if !isStrongPassword(password) {
		return errors.New("password must be at least 8 characters and include upper, lower, digit, and symbol")
	}

	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	if s.users.Exists(username) {
		return errors.New("username already exists")
	}

	if _, found := s.users.FindByEmail(email); found {
		return errors.New("email already registered")
	}

	return s.users.Add(&User{
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password),
		Profile: &UserProfile{
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
	if !exists || !verifyPassword(password, user.PasswordHash) {
		s.recordFailedLogin(username)
		return "", errors.New("invalid username or password")
	}

	s.clearFailedLogin(username)

	sessionID := generateSessionID(username)
	s.sessions.Create(&Session{
		Username:  username,
		SessionID: sessionID,
		IsActive:  true,
		ExpiresAt: time.Now().Add(30 * time.Minute),
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

	if time.Now().After(session.ExpiresAt) {
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
	return time.Now().Before(failure.LockedUntil)
}

func (s *authService) recordFailedLogin(username string) {
	failure := s.failures[username]
	now := time.Now()
	if now.Sub(failure.LastAttempt) > 15*time.Minute {
		failure.Count = 0
	}
	failure.Count++
	failure.LastAttempt = now
	if failure.Count >= 5 {
		failure.LockedUntil = now.Add(15 * time.Minute)
	}
	s.failures[username] = failure
}

func (s *authService) clearFailedLogin(username string) {
	delete(s.failures, username)
}

func (s *authService) GetUserInfo(sessionID string) (*User, error) {
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

	if !verifyPassword(oldPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if !isStrongPassword(newPassword) {
		return errors.New("new password must be at least 8 characters and include upper, lower, digit, and symbol")
	}

	user.PasswordHash = hashPassword(newPassword)
	return nil
}

func (s *authService) CreatePasswordResetToken(email string) (string, error) {
	user, exists := s.users.FindByEmail(email)
	if !exists {
		return "", errors.New("email not found")
	}
	token := generateSessionID(user.Username)
	s.resetTokens[token] = passwordResetToken{Username: user.Username, ExpiresAt: time.Now().Add(1 * time.Hour)}
	return token, nil
}

func (s *authService) ResetPasswordWithToken(token, newPassword string) error {
	reset, exists := s.resetTokens[token]
	if !exists {
		return errors.New("invalid token")
	}

	if time.Now().After(reset.ExpiresAt) {
		delete(s.resetTokens, token)
		return errors.New("reset token expired")
	}

	if !isStrongPassword(newPassword) {
		return errors.New("new password must be at least 8 characters and include upper, lower, digit, and symbol")
	}

	user, exists := s.users.Get(reset.Username)
	if !exists {
		return errors.New("user not found")
	}

	user.PasswordHash = hashPassword(newPassword)
	delete(s.resetTokens, token)
	return nil
}

func (s *authService) ListUsernames() []string {
	return s.users.List()
}

func (s *authService) DeleteUser(username string) error {
	return s.users.Delete(username)
}
