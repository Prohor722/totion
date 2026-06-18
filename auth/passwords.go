package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *authService) recordPasswordResetRequest(email string) error {
	now := s.clock.Now()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	request, exists := s.resetRequests[email]
	if !exists || now.After(request.WindowEnds) {
		s.resetRequests[email] = passwordResetRequest{Count: 1, WindowEnds: now.Add(util.PasswordResetRequestWindow)}
		return nil
	}

	if request.Count >= util.PasswordResetRequestLimit {
		return ErrTooManyResetRequests
	}

	request.Count++
	s.resetRequests[email] = request
	return nil
}

func (s *authService) CreatePasswordResetToken(email string) (string, error) {
	email = normalizeEmail(email)
	if email == "" {
		return "", ErrInvalidEmail
	}
	if !util.IsValidEmail(email) {
		return "", ErrInvalidEmail
	}

	if err := s.recordPasswordResetRequest(email); err != nil {
		return "", err
	}

	user, exists := s.users.FindByEmail(email)
	if !exists {
		return "", ErrEmailNotFound
	}

	token := util.GenerateSessionID(user.Username)
	hash := hashToken(token)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.resetTokens[hash] = passwordResetToken{Username: user.Username, ExpiresAt: s.clock.Now().Add(util.PasswordResetTokenTTL)}

	return token, nil
}

func (s *authService) ResetPasswordWithToken(token, newPassword string) error {
	hash := hashToken(token)

	s.mutex.RLock()
	reset, exists := s.resetTokens[hash]
	s.mutex.RUnlock()
	if !exists {
		return ErrInvalidToken
	}

	if s.clock.Now().After(reset.ExpiresAt) {
		s.mutex.Lock()
		delete(s.resetTokens, hash)
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

	s.mutex.Lock()
	delete(s.resetTokens, hash)
	s.mutex.Unlock()
	return nil
}
