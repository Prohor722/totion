package main

import "fmt"

func main() {
	testProgram()
	fmt.Println("Welcome to totion app")
	fmt.Println("--- Login System Demo ---")
	fmt.Println()

	// Register users
	fmt.Println("1. Registering users...")
	if err := DefaultAuth.Register("john_doe", "john@example.com", "password123"); err != nil {
		fmt.Println("Error:", err)
	}
	if err := DefaultAuth.Register("jane_smith", "jane@example.com", "securePass456"); err != nil {
		fmt.Println("Error:", err)
	}

	// Try registering with invalid credentials
	fmt.Println(DefaultAuth.Register("ab", "invalid@email.com", "short"))

	// Login attempt
	fmt.Println("2. Login attempt...")
	sessionID, err := DefaultAuth.Login("john_doe", "password123")
	if err != nil {
		fmt.Println("Login error:", err)
	} else {
		fmt.Printf("Session ID: %s\n", sessionID)

		// Validate session
		fmt.Println("3. Validating session...")
		if username, err := DefaultAuth.ValidateSession(sessionID); err == nil {
			fmt.Printf("✓ Session is valid for user: %s\n", username)
		}

		// Get user info
		fmt.Println("4. Retrieving user info...")
		user, err := DefaultAuth.GetUserInfo(sessionID)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
		}

		// Change password
		fmt.Println("5. Changing password...")
		if err := DefaultAuth.ChangePassword(sessionID, "password123", "newPassword789"); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✓ Password changed successfully")
		}

		// Logout
		fmt.Println("6. Logging out...")
		if err := DefaultAuth.Logout(sessionID); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✓ Logged out successfully")
		}

		// Verify session is invalid after logout
		fmt.Println("7. Verifying session after logout...")
		if _, err := DefaultAuth.ValidateSession(sessionID); err != nil {
			fmt.Println("✓ Session correctly invalidated")
		}
	}

	// List all users
	fmt.Println("8. All registered users:")
	for _, username := range DefaultAuth.ListUsernames() {
		fmt.Printf("  - %s\n", username)
	}
}
