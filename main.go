package main

import "fmt"

func main() {
	testProgram()
	fmt.Println("Welcome to totion app")
	fmt.Println("--- Login System Demo ---")
	fmt.Println()

	// Register users via AuthManager abstraction
	fmt.Println("1. Registering users...")
	if msg := authManager.RegisterUser("john_doe", "john@example.com", "password123"); msg != "" {
		fmt.Println(msg)
	}
	if msg := authManager.RegisterUser("jane_smith", "jane@example.com", "securePass456"); msg != "" {
		fmt.Println(msg)
	}

	// Try registering with invalid credentials
	fmt.Println(authManager.RegisterUser("ab", "invalid@email.com", "short"))

	// Login attempt
	fmt.Println("2. Login attempt...")
	sessionID, errMsg := authManager.LoginUser("john_doe", "password123")
	if errMsg != "" {
		fmt.Println("Login error:", errMsg)
	} else {
		fmt.Printf("Session ID: %s\n", sessionID)

		// Validate session
		fmt.Println("3. Validating session...")
		if valid, username := userManager.ValidateSession(sessionID); valid {
			fmt.Printf("✓ Session is valid for user: %s\n", username)
		}

		// Get user info
		fmt.Println("4. Retrieving user info...")
		user, msg := userManager.GetUserInfo(sessionID)
		if msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
		}

		// Change password
		fmt.Println("5. Changing password...")
		if msg := userManager.ChangePassword(sessionID, "password123", "newPassword789"); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Println("✓ Password changed successfully")
		}

		// Logout
		fmt.Println("6. Logging out...")
		if msg := authManager.LogoutUser(sessionID); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Println("✓ Logged out successfully")
		}

		// Verify session is invalid after logout
		fmt.Println("7. Verifying session after logout...")
		if valid, _ := userManager.ValidateSession(sessionID); !valid {
			fmt.Println("✓ Session correctly invalidated")
		}
	}

	// List all users
	fmt.Println("8. All registered users:")
	for _, username := range userManager.ListAllUsers() {
		fmt.Printf("  - %s\n", username)
	}
}
