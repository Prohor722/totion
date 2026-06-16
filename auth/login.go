package auth

import (
	"fmt"
	"strings"
	"sync"
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
	mutex    sync.RWMutex
	sessions map[string]*model.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]*model.Session)}
}

func (r *InMemorySessionRepository) Create(session *model.Session) {
	if session == nil {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sessions[session.SessionID] = session
}

func (r *InMemorySessionRepository) Get(sessionID string) (*model.Session, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	session, ok := r.sessions[sessionID]
	return session, ok
}

func (r *InMemorySessionRepository) Remove(sessionID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
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
	mutex       sync.RWMutex
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
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" || email == "" || password == "" {
		return ErrAllFieldsRequired
	}

	if strings.ContainsAny(username, " \t\n\r") {
		return ErrUsernameContainsWhitespace
	}

	if len(username) < 3 {
		return ErrUsernameTooShort
	}

	if !util.IsStrongPassword(password) {
		return ErrWeakPassword
	}

	if !util.IsValidEmail(email) {
		return ErrInvalidEmail
	}

	if s.users.Exists(username) {
		return ErrUsernameAlreadyExists
	}

	if _, found := s.users.FindByEmail(email); found {
		return ErrEmailAlreadyRegistered
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
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", ErrAllFieldsRequired
	}

	if s.isAccountLocked(username) {
		return "", ErrAccountLocked
	}

	user, exists := s.users.Get(username)
	if !exists || !util.VerifyPassword(password, user.PasswordHash) {
		s.recordFailedLogin(username)
		return "", ErrInvalidCredentials
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
	if sessionID == "" {
		return ErrSessionIDRequired
	}
	s.sessions.Remove(sessionID)
	return nil
}

func (s *authService) ValidateSession(sessionID string) (string, error) {
	if sessionID == "" {
		return "", ErrSessionIDRequired
	}

	session, exists := s.sessions.Get(sessionID)
	if !exists || session == nil || !session.IsActive {
		return "", ErrInvalidOrExpiredSession
	}

	if s.clock.Now().After(session.ExpiresAt) {
		s.sessions.Remove(sessionID)
		return "", ErrSessionExpired
	}

	return session.Username, nil
}

func (s *authService) isAccountLocked(username string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	failure, exists := s.failures[username]
	if !exists {
		return false
	}
	return s.clock.Now().Before(failure.LockedUntil)
}

func (s *authService) recordFailedLogin(username string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.failures, username)
}

func (s *authService) GetUserInfo(sessionID string) (*model.User, error) {
	username, err := s.ValidateSession(sessionID)
	if err != nil {
		return nil, err
	}

	user, exists := s.users.Get(username)
	if !exists {
		return nil, ErrUserNotFound
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
		return ErrUserNotFound
	}

	if !util.VerifyPassword(oldPassword, user.PasswordHash) {
		return ErrCurrentPasswordIncorrect
	}

	if !util.IsStrongPassword(newPassword) {
		return ErrWeakPassword
	}

	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	return nil
}

func (s *authService) CreatePasswordResetToken(email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrInvalidEmail
	}

	user, exists := s.users.FindByEmail(email)
	if !exists {
		return "", ErrEmailNotFound
	}
	token := util.GenerateSessionID(user.Username)

	s.mutex.Lock()
	s.resetTokens[token] = passwordResetToken{Username: user.Username, ExpiresAt: s.clock.Now().Add(util.PasswordResetTokenTTL)}
	s.mutex.Unlock()

	return token, nil
}

func (s *authService) ResetPasswordWithToken(token, newPassword string) error {
	s.mutex.RLock()
	reset, exists := s.resetTokens[token]
	s.mutex.RUnlock()
	if !exists {
		return ErrInvalidToken
	}

	if s.clock.Now().After(reset.ExpiresAt) {
		s.mutex.Lock()
		delete(s.resetTokens, token)
		s.mutex.Unlock()
		return ErrResetTokenExpired
	}

	if !util.IsStrongPassword(newPassword) {
		return ErrWeakPassword
	}

	user, exists := s.users.Get(reset.Username)
	if !exists {
		return ErrUserNotFound
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
