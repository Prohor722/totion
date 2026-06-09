package auth

import (
	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
	"github.com/Prohor722/totion/util"
)

type ProfileService interface {
	GetProfile(username string) (*model.UserProfile, error)
	UpdateProfile(username, email, bio string) error
}

// profileService implements ProfileService using a UserRepository.
type profileService struct {
	users store.UserRepository
}

func NewProfileService(users store.UserRepository) ProfileService {
	return &profileService{users: users}
}

func (s *profileService) GetProfile(username string) (*model.UserProfile, error) {
	user, exists := s.users.Get(username)
	if !exists {
		return nil, ErrUserNotFound
	}
	if user.Profile == nil {
		return nil, ErrInvalidInput
	}
	return &model.UserProfile{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Profile.Bio,
	}, nil
}

func (s *profileService) UpdateProfile(username, email, bio string) error {
	user, exists := s.users.Get(username)
	if !exists {
		return ErrUserNotFound
	}

	if email != "" {
		if !util.IsValidEmail(email) {
			return ErrInvalidInput
		}

		if existing, found := s.users.FindByEmail(email); found && existing.Username != username {
			return ErrUserExists
		}

		user.Email = email
		if user.Profile != nil {
			user.Profile.Email = email
		}
	}

	if bio != "" {
		if user.Profile == nil {
			user.Profile = &model.UserProfile{Username: user.Username, Email: user.Email}
		}
		user.Profile.Bio = bio
	}

	return nil
}

type ProfileManager interface {
	GetUserProfile(username string) (*model.UserProfile, error)
	UpdateUserProfile(username, email, bio string) error
}

type defaultProfileManager struct {
	service ProfileService
}

func NewDefaultProfileManager(service ProfileService) ProfileManager {
	return &defaultProfileManager{service: service}
}

func (m *defaultProfileManager) GetUserProfile(username string) (*model.UserProfile, error) {
	return m.service.GetProfile(username)
}

func (m *defaultProfileManager) UpdateUserProfile(username, email, bio string) error {
	return m.service.UpdateProfile(username, email, bio)
}
