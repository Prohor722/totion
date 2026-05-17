package main

import "errors"

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
	ListUsernames() []string
	DeleteUser(username string) error
}

type authService struct {
	users    UserRepository
	sessions SessionRepository
}

func NewAuthService(users UserRepository, sessions SessionRepository) AuthService {
	return &authService{users: users, sessions: sessions}
}

func (s *authService) Register(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("all fields are required")
	}

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
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

	user, exists := s.users.Get(username)
	if !exists || !verifyPassword(password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	sessionID := generateSessionID(username)
	s.sessions.Create(&Session{
		Username:  username,
		SessionID: sessionID,
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

	if !verifyPassword(oldPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	user.PasswordHash = hashPassword(newPassword)
	return nil
}

func (s *authService) ListUsernames() []string {
	return s.users.List()
}

func (s *authService) DeleteUser(username string) error {
	return s.users.Delete(username)
}

