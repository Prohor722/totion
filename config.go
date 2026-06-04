package main

import "time"

// Centralized configuration for auth and password policies.
var (
    MinPasswordLength     = 8
    SessionTTL            = 30 * time.Minute
    FailedLoginThreshold  = 5
    FailedLoginWindow     = 15 * time.Minute
    AccountLockDuration   = 15 * time.Minute
    PasswordResetTokenTTL = 1 * time.Hour
)
package main

import "time"

// Configuration constants for security and session behavior.
const (
	MinPasswordLength      = 8
	SessionTTL             = 30 * time.Minute
	PasswordResetTokenTTL  = 1 * time.Hour
	MaxFailedLoginAttempts = 5
	FailedLoginWindow      = 15 * time.Minute
	LockoutDuration        = 15 * time.Minute
)
