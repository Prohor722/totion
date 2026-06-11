package auth

import (
	"testing"
	"time"
	"github.com/Prohor722/totion/store"
)

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	return c.now
}

func (c *testClock) Advance(d time.Duration) {
	c.now = c.now.Add(d)
}

func TestAuthService_RegisterLoginLogoutValidate(t *testing.T) {
	clock := &testClock{now: time.Now()}
	users := store.NewInMemoryUserRepository()
	sessions := NewInMemorySessionRepository()
	auth := NewAuthServiceWithClock(users, sessions, clock)

	if err := auth.Register("alice", "alice@example.com", "Secure1!"); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	sessionID, err := auth.Login("alice", "Secure1!")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	username, err := auth.ValidateSession(sessionID)
	if err != nil {
		t.Fatalf("ValidateSession failed: %v", err)
	}
	if username != "alice" {
		t.Fatalf("expected username alice, got %s", username)
	}

	if err := auth.Logout(sessionID); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if _, err := auth.ValidateSession(sessionID); err == nil {
		t.Fatal("expected session validation to fail after logout")
	}
}

func TestAuthService_SessionExpires(t *testing.T) {
	clock := &testClock{now: time.Now()}
	users := store.NewInMemoryUserRepository()
	sessions := NewInMemorySessionRepository()
	auth := NewAuthServiceWithClock(users, sessions, clock)

	if err := auth.Register("bob", "bob@example.com", "Strong1$"); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	sessionID, err := auth.Login("bob", "Strong1$")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	clock.Advance(SessionTTL + time.Minute)
	if _, err := auth.ValidateSession(sessionID); err == nil {
		t.Fatal("expected session to expire")
	}
}

func TestAuthService_PasswordResetFlow(t *testing.T) {
	clock := &testClock{now: time.Now()}
	users := store.NewInMemoryUserRepository()
	sessions := NewInMemorySessionRepository()
	auth := NewAuthServiceWithClock(users, sessions, clock)

	if err := auth.Register("dave", "dave@example.com", "Reset1@A"); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	token, err := auth.CreatePasswordResetToken("dave@example.com")
	if err != nil {
		t.Fatalf("CreatePasswordResetToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected reset token")
	}

	if err := auth.ResetPasswordWithToken(token, "NewPass1#"); err != nil {
		t.Fatalf("ResetPasswordWithToken failed: %v", err)
	}

	if _, err := auth.Login("dave", "NewPass1#"); err != nil {
		t.Fatalf("Login with new password failed: %v", err)
	}
}

func TestAuthService_PasswordResetTokenExpires(t *testing.T) {
	clock := &testClock{now: time.Now()}
	users := store.NewInMemoryUserRepository()
	sessions := NewInMemorySessionRepository()
	auth := NewAuthServiceWithClock(users, sessions, clock)

	if err := auth.Register("evan", "evan@example.com", "Reset2#A"); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	token, err := auth.CreatePasswordResetToken("evan@example.com")
	if err != nil {
		t.Fatalf("CreatePasswordResetToken failed: %v", err)
	}

	clock.Advance(PasswordResetTokenTTL + time.Minute)
	if err := auth.ResetPasswordWithToken(token, "Another1$"); err == nil {
		t.Fatal("expected reset token to expire")
	}
}
