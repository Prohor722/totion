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

type AuthService interface {
	Register(username, email, password string) error
	Login(username, password string) (string, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (string, error)
	GetUserInfo(sessionID string) (*User, error)
	ChangePassword(sessionID, oldPassword, newPassword string) error
	CreatePasswordResetToken(email string) (string, error)
	ResetPasswordWithToken(token, newPassword string) error
	ListUsernames() []string
	DeleteUser(username string) error
}

const (
	sessionExpiryDuration     = 30 * time.Minute
	resetTokenExpiryDuration  = 15 * time.Minute
	maxFailedLoginAttempts    = 5
	accountLockoutDuration    = 15 * time.Minute
)

type resetTokenInfo struct {
	username  string
	expiresAt time.Time
}

type authService struct {
	users       UserRepository
	sessions    SessionRepository
	resetTokens map[string]resetTokenInfo
}

func NewAuthService(users UserRepository, sessions SessionRepository) AuthService {
	return &authService{users: users, sessions: sessions, resetTokens: make(map[string]resetTokenInfo)}
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

	if err := validatePasswordPolicy(password); err != nil {
		return err
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

	salt := generateSalt()
	return s.users.Add(&User{
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password, salt),
		PasswordSalt: salt,
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

	user, exists := s.users.Get(username)
	if !exists {
		return "", errors.New("invalid username or password")
	}

	if user.LockoutUntil.After(time.Now()) {
		return "", errors.New("account locked due to repeated failed login attempts")
	}

	if !verifyPassword(password, user.PasswordSalt, user.PasswordHash) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= maxFailedLoginAttempts {
			user.LockoutUntil = time.Now().Add(accountLockoutDuration)
			return "", errors.New("account locked due to repeated failed login attempts")
		}
		return "", errors.New("invalid username or password")
	}

	user.FailedLoginAttempts = 0
	user.LockoutUntil = time.Time{}

	sessionID := generateSessionID(username)
	s.sessions.Create(&Session{
		Username:  username,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(sessionExpiryDuration),
		IsActive:  true,
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
		sessions := s.sessions
		session.IsActive = false
		sessions.Remove(sessionID)
		return "", errors.New("session has expired")
	}

	return session.Username, nil
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

	if !verifyPassword(oldPassword, user.PasswordSalt, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if err := validatePasswordPolicy(newPassword); err != nil {
		return err
	}

	salt := generateSalt()
	user.PasswordSalt = salt
	user.PasswordHash = hashPassword(newPassword, salt)
	return nil
}

func (s *authService) CreatePasswordResetToken(email string) (string, error) {
	user, exists := s.users.FindByEmail(email)
	if !exists {
		return "", errors.New("email not found")
	}
	token := generateSecureToken()
	s.resetTokens[token] = resetTokenInfo{
		username:  user.Username,
		expiresAt: time.Now().Add(resetTokenExpiryDuration),
	}
	return token, nil
}

func (s *authService) ResetPasswordWithToken(token, newPassword string) error {
	info, exists := s.resetTokens[token]
	if !exists {
		return errors.New("invalid token")
	}

	if time.Now().After(info.expiresAt) {
		delete(s.resetTokens, token)
		return errors.New("token has expired")
	}

	if err := validatePasswordPolicy(newPassword); err != nil {
		return err
	}

	user, exists := s.users.Get(info.username)
	if !exists {
		return errors.New("user not found")
	}

	salt := generateSalt()
	user.PasswordSalt = salt
	user.PasswordHash = hashPassword(newPassword, salt)
	delete(s.resetTokens, token)
	return nil
}

func (s *authService) ListUsernames() []string {
	return s.users.List()
}

func (s *authService) DeleteUser(username string) error {
	return s.users.Delete(username)
}
