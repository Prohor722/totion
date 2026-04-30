package main

// GetUserInfo retrieves user information (requires active session)
func GetUserInfo(sessionID string) (*User, string) {
	isValid, username := ValidateSession(sessionID)
	if !isValid {
		return nil, "Error: Invalid or expired session"
	}

	user, exists := userDatabase[username]
	if !exists {
		return nil, "Error: User not found"
	}

	return user, ""
}

// ListAllUsers returns all registered users (for admin purposes)
func ListAllUsers() []string {
	var usernames []string
	for username := range userDatabase {
		usernames = append(usernames, username)
	}
	return usernames
}