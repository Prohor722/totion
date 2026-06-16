package auth

import (
    "strings"

    "github.com/Prohor722/totion/model"
    "github.com/Prohor722/totion/util"
)

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

func (s *authService) ListUsernames() []string {
    return s.users.List()
}

func (s *authService) DeleteUser(username string) error {
    return s.users.Delete(username)
}
