package util

import (
	"encoding/hex"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	t.Parallel()

	validCases := []string{
		"user@example.com",
		"user+tag@example.com",
		"firstname.lastname@example.co.uk",
	}

	for _, email := range validCases {
		email := email
		t.Run(email, func(t *testing.T) {
			t.Parallel()
			if !IsValidEmail(email) {
				t.Fatalf("expected valid email for %s", email)
			}
		})
	}

	invalidCases := []string{
		"userexample.com",
		"user@ example.com",
		"user@com",
		"@example.com",
		"user@.com",
		" user@example.com ",
	}

	for _, email := range invalidCases {
		email := email
		t.Run(email, func(t *testing.T) {
			t.Parallel()
			if IsValidEmail(email) {
				t.Fatalf("expected invalid email for %s", email)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	t.Parallel()

	valid := "Abcdef1!"
	if !IsStrongPassword(valid) {
		t.Fatalf("expected password to be strong: %s", valid)
	}

	invalidCases := []struct {
		name     string
		password string
	}{
		{"too short", "Ab1!"},
		{"missing upper", "abcdef1!"},
		{"missing lower", "ABCDEF1!"},
		{"missing digit", "Abcdefg!"},
		{"missing symbol", "Abcdef12"},
		{"contains whitespace", "Abc 1!ef"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if IsStrongPassword(tc.password) {
				t.Fatalf("expected password to be weak: %s", tc.password)
			}
		})
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	t.Parallel()

	password := "Abcdef1!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !VerifyPassword(password, hash) {
		t.Fatal("expected VerifyPassword to validate the correct password")
	}

	if VerifyPassword("wrong", hash) {
		t.Fatal("expected VerifyPassword to reject an incorrect password")
	}
}

func TestGenerateSessionID(t *testing.T) {
	t.Parallel()

	ids := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		id := GenerateSessionID("user")
		if id == "" {
			t.Fatal("expected non-empty session ID")
		}
		if len(id) != sessionIDSize*2 {
			t.Fatalf("expected session ID length %d, got %d", sessionIDSize*2, len(id))
		}
		if _, err := hex.DecodeString(id); err != nil {
			t.Fatalf("expected session ID to be valid hex, got %q: %v", id, err)
		}
		if _, exists := ids[id]; exists {
			t.Fatal("expected generated session IDs to be unique")
		}
		ids[id] = struct{}{}
	}
}
