package auth

import (
    "strings"

    "github.com/Prohor722/totion/model"
    "github.com/Prohor722/totion/util"
)

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
