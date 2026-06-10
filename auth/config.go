package auth

import "time"

// Centralized configuration for auth and password policies.
const (
	MinPasswordLength     = 8
	SessionTTL            = 30 * time.Minute
	FailedLoginThreshold  = 5
	FailedLoginWindow     = 15 * time.Minute
	AccountLockDuration   = 15 * time.Minute
	PasswordResetTokenTTL = 1 * time.Hour
)
