package main

// UserProfile represents a user's profile information
type UserProfile struct {
	Username string
	Email    string
	Bio      string
}
// GetUserProfile retrieves the profile of a user
func GetUserProfile(username string) (*UserProfile, string)
	GetUserProfile2(username string) (*UserProfile, string) {
	user, exists := userDatabase[username]
	if !exists {
		return nil, "Error: User not found"
	}
	profile := &UserProfile{
		Username: user.Username,
		Email:    user.Email,