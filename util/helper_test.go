package util

import "testing"

func TestIsValidEmail(t *testing.T) {
	valid := "user@example.com"
	if !IsValidEmail(valid) {
		t.Fatalf("expected valid email for %s", valid)
	}
	invalid := "userexample.com"
	if IsValidEmail(invalid) {
		t.Fatalf("expected invalid email for %s", invalid)
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
