package main

import "time"

// Configuration constants for security and session behavior.
const (
	MinPasswordLength = 8
	SessionTTL        = 30 * time.Minute
)
