package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const sessionIDSize = 16

var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(password, hash string) bool {
	if hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func IsStrongPassword(password string) bool {
	password = strings.TrimSpace(password)
	if len(password) < MinPasswordLength {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, ch := range password {
		switch {
		case unicode.IsSpace(ch):
			return false
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSymbol = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSymbol
}

func GenerateSessionID(username string) string {
	buffer := make([]byte, sessionIDSize)
	if n, err := rand.Read(buffer); err == nil && n == len(buffer) {
		return hex.EncodeToString(buffer)
	}

	fallbackSource := fmt.Sprintf("%s:%d", username, time.Now().UnixNano())
	fallback := sha256.Sum256([]byte(fallbackSource))
	return hex.EncodeToString(fallback[:])[:sessionIDSize*2]
}

func IsValidEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func IsValidWebsite(website string) bool {
	website = strings.TrimSpace(website)
	if website == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(website)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
