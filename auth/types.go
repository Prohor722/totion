package auth

import "time"

// Clock provides time abstraction for easier testing.
type Clock interface {
    Now() time.Time
}

// passwordResetToken holds temporary token metadata for password resets.
type passwordResetToken struct {
    Username  string
    ExpiresAt time.Time
}

// failedLogin tracks failed attempts and lockout windows for a user.
type failedLogin struct {
    Count       int
    LastAttempt time.Time
    LockedUntil time.Time
}
