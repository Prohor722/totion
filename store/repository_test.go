package store

import (
	"testing"
	"github.com/Prohor722/totion/model"
)

func TestInMemoryUserRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := &User{Username: "alice", Email: "alice@example.com", PasswordHash: "h"}
	if err := repo.Add(u); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if !repo.Exists("alice") {
		t.Fatalf("Exists should return true after add")
	}

	got, ok := repo.Get("alice")
	if !ok || got.Username != "alice" {
		t.Fatalf("Get returned wrong user: %v, %v", got, ok)
	}

	byEmail, found := repo.FindByEmail("alice@example.com")
	if !found || byEmail.Username != "alice" {
		t.Fatalf("FindByEmail failed: %v, %v", byEmail, found)
	}

	list := repo.List()
	if len(list) != 1 || list[0] != "alice" {
		t.Fatalf("List unexpected: %v", list)
	}

	if err := repo.Delete("alice"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if repo.Exists("alice") {
		t.Fatalf("Exists should be false after delete")
	}
}
