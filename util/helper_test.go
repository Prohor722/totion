package util

import "testing"

func TestIsValidEmail(t *testing.T) {
	valid := "user@example.com"
	if !IsValidEmail(valid) {
		t.Fatalf("expected valid email for %s", valid)
	}
	plusAddr := "user+tag@example.com"
	if !IsValidEmail(plusAddr) {
		t.Fatalf("expected valid email for %s", plusAddr)
	}

	invalids := []string{"userexample.com", "user@ example.com", "user@com", "@example.com"}
	for _, input := range invalids {
		if IsValidEmail(input) {
			t.Fatalf("expected invalid email for %s", input)
		}
	}
}

func TestIsStrongPasswordAndHashVerify(t *testing.T) {
	pass := "Abcdef1!"
	if !IsStrongPassword(pass) {
		t.Fatalf("expected password to be strong: %s", pass)
	}

	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !VerifyPassword(pass, hash) {
		t.Fatalf("VerifyPassword failed for valid password/hash")
	}

	if VerifyPassword("wrong", hash) {
		t.Fatalf("VerifyPassword succeeded for wrong password")
	}
}

func TestIsStrongPasswordRejectsWhitespace(t *testing.T) {
	if IsStrongPassword("Abc 1!ef") {
		t.Fatal("expected password with spaces to be rejected")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1 := GenerateSessionID("alice")
	id2 := GenerateSessionID("bob")

	if id1 == "" || id2 == "" {
		t.Fatal("expected non-empty session IDs")
	}
	if id1 == id2 {
		t.Fatal("expected generated session IDs to be unique")
	}
	if len(id1) != sessionIDSize*2 {
		t.Fatalf("expected session ID length %d, got %d", sessionIDSize*2, len(id1))
	}
}
