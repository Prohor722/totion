package main

import "time"

// Session represents a user session
type Session struct {
	Username  string
	SessionID string
	IsActive  bool
	ExpiresAt time.Time
}
