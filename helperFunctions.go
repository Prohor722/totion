package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

func hashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(salt + password))
	return fmt.Sprintf("%x", hash)
}

func verifyPassword(password, salt, hash string) bool {
	return hashPassword(password, salt) == hash
}

func generateSalt() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err == nil {
		return hex.EncodeToString(buffer)
	}

	fallback := sha256.Sum256([]byte(fmt.Sprintf("fallback-%d", len(buffer))))
	return hex.EncodeToString(fallback[:])[:32]
}

func generateSecureToken() string {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err == nil {
		return hex.EncodeToString(buffer)
	}

	fallback := sha256.Sum256([]byte(fmt.Sprintf("fallback-token-%d", len(buffer))))
	return hex.EncodeToString(fallback[:])
}

func generateSessionID(username string) string {
	return generateSecureToken()
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validatePasswordPolicy(password string) error {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must include uppercase, lowercase, digit, and special character")
	}
	return nil
}
