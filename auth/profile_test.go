package auth

import (
	"testing"

	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
)

func TestProfileService_UpdateWebsite(t *testing.T) {
	users := store.NewInMemoryUserRepository()
	service := NewProfileService(users)

	user := &model.User{
		Username: "gina",
		Email:    "gina@example.com",
		Profile: &model.UserProfile{
			Username: "gina",
			Email:    "gina@example.com",
		},
	}
	if err := users.Add(user); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	website := "https://gina.example.com"
	if err := service.UpdateProfile("gina", ProfileUpdate{Website: &website}); err != nil {
		t.Fatalf("expected UpdateProfile to succeed, got %v", err)
	}

	profile, err := service.GetProfile("gina")
	if err != nil {
		t.Fatalf("expected GetProfile to succeed, got %v", err)
	}
	if profile.Website != website {
		t.Fatalf("expected website %q, got %q", website, profile.Website)
	}
}

func TestProfileService_UpdateInvalidWebsite(t *testing.T) {
	users := store.NewInMemoryUserRepository()
	service := NewProfileService(users)

	user := &model.User{
		Username: "harry",
		Email:    "harry@example.com",
		Profile: &model.UserProfile{
			Username: "harry",
			Email:    "harry@example.com",
		},
	}
	if err := users.Add(user); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	website := "ftp://invalid.example.com"
	if err := service.UpdateProfile("harry", ProfileUpdate{Website: &website}); err == nil {
		t.Fatal("expected invalid website update to fail")
	}
}
