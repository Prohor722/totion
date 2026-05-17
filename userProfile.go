package main

import "errors"

type User struct {
	Username     string
	Email        string
	PasswordHash string
	Profile      *UserProfile
}

type UserProfile struct {
	Username string
	Email    string
	Bio      string
}

type ProfileService interface {
	GetProfile(username string) (*UserProfile, error)
	UpdateProfile(username, email, bio string) error
}

type profileService struct {
	users UserRepository
}

func NewProfileService(users UserRepository) ProfileService {
	return &profileService{users: users}
}

func (s *profileService) GetProfile(username string) (*UserProfile, error) {
	user, exists := s.users.Get(username)
	if !exists {
		return nil, errors.New("user not found")
	}
	return &UserProfile{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Profile.Bio,
	}, nil
}

func (s *profileService) UpdateProfile(username, email, bio string) error {
	user, exists := s.users.Get(username)
	if !exists {
		return errors.New("user not found")
	}

	if email != "" {
		if !isValidEmail(email) {
			return errors.New("invalid email format")
		}

		if existing, found := s.users.FindByEmail(email); found && existing.Username != username {
			return errors.New("email already registered")
		}

		user.Email = email

		user.Profile.Email = email
	}

	if bio != "" {
		user.Profile.Bio = bio
	}

	return nil
}
