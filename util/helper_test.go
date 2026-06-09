package main

import "testing"

func TestIsValidEmail(t *testing.T) {
	valid := "user@example.com"
	if !isValidEmail(valid) {
		t.Fatalf("expected valid email for %s", valid)
	}
	invalid := "userexample.com"
	if isValidEmail(invalid) {
		t.Fatalf("expected invalid email for %s", invalid)
	}
}

func TestIsStrongPasswordAndHashVerify(t *testing.T) {
	pass := "Abcdef1!"
	if !isStrongPassword(pass) {
		t.Fatalf("expected password to be strong: %s", pass)
	}

	hash, err := hashPassword(pass)
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}

	if !verifyPassword(pass, hash) {
		t.Fatalf("verifyPassword failed for valid password/hash")
	}

	if verifyPassword("wrong", hash) {
		t.Fatalf("verifyPassword succeeded for wrong password")
	}
}
