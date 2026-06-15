package auth

import (
	"errors"
	"fmt"

	"github.com/Prohor722/totion/util"
)

var (
	ErrUserExists                 = errors.New("username already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrUserIsNil                  = errors.New("user is nil")
	ErrInvalidInput               = errors.New("invalid input")
	ErrUsernameContainsWhitespace = errors.New("username cannot contain whitespace")

	ErrAllFieldsRequired        = errors.New("all fields are required")
	ErrUsernameTooShort         = errors.New("username must be at least 3 characters")
	ErrWeakPassword             = fmt.Errorf("password must be at least %d characters and include upper, lower, digit, and symbol", util.MinPasswordLength)
	ErrInvalidEmail             = errors.New("invalid email format")
	ErrUsernameAlreadyExists    = ErrUserExists
	ErrEmailAlreadyRegistered   = errors.New("email already registered")
	ErrInvalidCredentials       = errors.New("invalid username or password")
	ErrAccountLocked            = errors.New("account temporarily locked due to failed login attempts")
	ErrSessionNotFound          = errors.New("session not found")
	ErrInvalidOrExpiredSession  = errors.New("invalid or expired session")
	ErrSessionExpired           = errors.New("session expired")
	ErrCurrentPasswordIncorrect = errors.New("current password is incorrect")
	ErrInvalidToken             = errors.New("invalid token")
	ErrResetTokenExpired        = errors.New("reset token expired")
	ErrEmailNotFound            = errors.New("email not found")
	ErrSessionIDRequired        = errors.New("session ID is required")
)
