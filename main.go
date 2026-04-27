package main

import "fmt"

func main() {
	testProgram()
	fmt.Println("Welcome to totion app")
	fmt.Println("\n--- Login System Demo ---\n")

	// Register users
	fmt.Println("1. Registering users...")
	RegisterUser("john_doe", "john@example.com", "password123")
	RegisterUser("jane_smith", "jane@example.com", "securePass456")

	// Try registering with invalid credentials
	fmt.Println(RegisterUser("ab", "invalid@email.com", "short"))

	// Login attempt
	fmt.Println("\n2. Login attempt...")
	sessionID, err := LoginUser("john_doe", "password123")
	if err == "" {
		fmt.Printf("Session ID: %s\n", sessionID)

		// Validate session
		fmt.Println("\n3. Validating session...")
		if isValid, username := ValidateSession(sessionID); isValid {
			fmt.Printf("✓ Session is valid for user: %s\n", username)
		}

		// Get user info
		fmt.Println("\n4. Retrieving user info...")
		user, _ := GetUserInfo(sessionID)
		if user != nil {
			fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
		}

		// Change password
		fmt.Println("\n5. Changing password...")
		fmt.Println(ChangePassword(sessionID, "password123", "newPassword789"))

		// Logout
		fmt.Println("\n6. Logging out...")
		fmt.Println(LogoutUser(sessionID))

		// Verify session is invalid after logout
		fmt.Println("\n7. Verifying session after logout...")
		if isValid, _ := ValidateSession(sessionID); !isValid {
			fmt.Println("✓ Session correctly invalidated")
		}
	}

	// List all users
	fmt.Println("\n8. All registered users:")
	for _, username := range ListAllUsers() {
		fmt.Printf("  - %s\n", username)
	}
}