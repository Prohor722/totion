// RegisterUser creates a new user account
// Returns error message if registration fails, empty string if successful
func RegisterUser(username, email, password string) string {
	// Validation checks
	if username == "" || email == "" || password == "" {
		return "Error: All fields are required"
	}

	if len(username) < 3 {
		return "Error: Username must be at least 3 characters"
	}

	if len(password) < 6 {
		return "Error: Password must be at least 6 characters"
	}

	if !strings.Contains(email, "@") {
		return "Error: Invalid email format"
	}

	// Check if username already exists
	if _, exists := userDatabase[username]; exists {
		return "Error: Username already exists"
	}

	// Check if email already exists
	for _, user := range userDatabase {
		if user.Email == email {
			return "Error: Email already registered"
		}
	}

	// Hash the password
	hashedPassword := hashPassword(password)

	// Create and store the user
	newUser := &User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	userDatabase[username] = newUser
	fmt.Printf("✓ User '%s' registered successfully\n", username)
	return ""
}