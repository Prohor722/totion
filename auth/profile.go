package auth

import (
	"strings"

	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
	"github.com/Prohor722/totion/util"
)

type ProfileUpdate struct {
	Email   *string
	Bio     *string
	Website *string
}

type ProfileService interface {
	GetProfile(username string) (*model.UserProfile, error)
	UpdateProfile(username string, update ProfileUpdate) error
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
		Website:  user.Profile.Website,
	}, nil
}

func (s *profileService) ensureProfile(user *model.User) *model.UserProfile {
	if user.Profile == nil {
		user.Profile = &model.UserProfile{Username: user.Username, Email: user.Email}
	}
	return user.Profile
}

func (s *profileService) UpdateProfile(username string, update ProfileUpdate) error {
	user, exists := s.users.Get(username)
	if !exists {
		return ErrUserNotFound
	}

	if update.Email != nil {
		email := strings.TrimSpace(*update.Email)
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

	if update.Bio != nil {
		bio := strings.TrimSpace(*update.Bio)
		profile := s.ensureProfile(user)
		profile.Bio = bio
	}

	if update.Website != nil {
		website := strings.TrimSpace(*update.Website)
		if website != "" && !util.IsValidWebsite(website) {
			return ErrInvalidInput
		}
		profile := s.ensureProfile(user)
		profile.Website = website
	}

	return nil
}

type ProfileManager interface {
	GetUserProfile(username string) (*model.UserProfile, error)
	UpdateUserProfile(username string, update ProfileUpdate) error
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

func (m *defaultProfileManager) UpdateUserProfile(username string, update ProfileUpdate) error {
	return m.service.UpdateProfile(username, update)
}
