package main

import (
	"crypto/sha256"
	"fmt"
)

// hashPassword creates a SHA256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// verifyPassword checks if the provided password matches the hash
func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}