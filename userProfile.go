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

// ProfileService defines low-level profile data operations.
// Focused on persistence and retrieval - Single Responsibility Principle.
type ProfileService interface {
	GetProfile(username string) (*UserProfile, error)
	UpdateProfile(username, email, bio string) error
}

// profileService implements ProfileService using a UserRepository.
// Dependency Inversion: depends on UserRepository abstraction, not concrete implementations.
type profileService struct {
	users UserRepository
}

func NewProfileService(users UserRepository) ProfileService {
	return &profileService{users: users}
}

// DefaultProfileService is the package-level profile service.
var DefaultProfileService ProfileService = NewProfileService(UserStore)

func (s *profileService) GetProfile(username string) (*UserProfile, error) {
	user, exists := s.users.Get(username)
	if !exists {
		return nil, errors.New("user not found")
	}
	if user.Profile == nil {
		return nil, errors.New("profile data unavailable")
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

// ProfileManager defines high-level profile management operations.
// This interface follows Interface Segregation Principle - clients depend on specific operations.
type ProfileManager interface {
	GetUserProfile(username string) (*UserProfile, error)
	UpdateUserProfile(username, email, bio string) error
}

// defaultProfileManager implements ProfileManager by delegating to ProfileService.
// This follows Dependency Inversion Principle - concrete implementation depends on abstraction.
type defaultProfileManager struct {
	service ProfileService
}

// NewDefaultProfileManager creates a new ProfileManager with the given ProfileService.
func NewDefaultProfileManager(service ProfileService) ProfileManager {
	return &defaultProfileManager{service: service}
}

// GetUserProfile retrieves a user's profile by username.
// This method is responsible for profile retrieval only - Single Responsibility Principle.
func (m *defaultProfileManager) GetUserProfile(username string) (*UserProfile, error) {
	return m.service.GetProfile(username)
}

// UpdateUserProfile updates a user's profile information.
// This method is responsible for profile updates only - Single Responsibility Principle.
func (m *defaultProfileManager) UpdateUserProfile(username, email, bio string) error {
	return m.service.UpdateProfile(username, email, bio)
}

// Legacy package-level functions that use the default global ProfileManager.
// These maintain backwards compatibility with existing code.

// profileManager is the default global instance.
var profileManager ProfileManager = NewDefaultProfileManager(DefaultProfileService)

// GetUserProfile retrieves a user's profile by username.
func GetUserProfile(username string) (*UserProfile, error) {
	return profileManager.GetUserProfile(username)
}

// UpdateUserProfile updates a user's profile information.
func UpdateUserProfile(username, email, bio string) error {
	return profileManager.UpdateUserProfile(username, email, bio)
}
