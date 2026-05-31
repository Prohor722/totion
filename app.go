package main

import "fmt"

type App struct {
	auth           AuthManager
	userManager    UserManager
	profileManager ProfileManager
	resetService   ForgetPasswordService
}

func NewApp(authService AuthService, userRepo UserRepository) *App {
	return &App{
		auth:           NewDefaultAuthManager(authService),
		userManager:    NewDefaultUserManager(authService),
		profileManager: NewDefaultProfileManager(NewProfileService(userRepo)),
		resetService:   NewForgetPasswordService(authService),
	}
}

func (a *App) RunDemo() {
	fmt.Println("Welcome to totion app")
	fmt.Println("--- Login System Demo ---")
	fmt.Println()

	fmt.Println("1. Registering users...")
	if err := a.auth.RegisterUser("john_doe", "john@example.com", "password123"); err != nil {
		fmt.Println("Register error:", err)
	}
	if err := a.auth.RegisterUser("jane_smith", "jane@example.com", "securePass456"); err != nil {
		fmt.Println("Register error:", err)
	}

	if err := a.auth.RegisterUser("ab", "invalid@email.com", "short"); err != nil {
		fmt.Println("Register error:", err)
	}

	fmt.Println("2. Login attempt...")
	sessionID, err := a.auth.LoginUser("john_doe", "password123")
	if err != nil {
		fmt.Println("Login error:", err)
	} else {
		fmt.Printf("Session ID: %s\n", sessionID)

		fmt.Println("3. Validating session...")
		username, err := a.userManager.ValidateSession(sessionID)
		if err == nil {
			fmt.Printf("✓ Session is valid for user: %s\n", username)
		}

		fmt.Println("4. Retrieving user info...")
		user, err := a.userManager.GetUserInfo(sessionID)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
		}

		fmt.Println("5. Changing password...")
		if err := a.userManager.ChangePassword(sessionID, "password123", "newPassword789"); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✓ Password changed successfully")
		}

		fmt.Println("6. Logging out...")
		if err := a.auth.LogoutUser(sessionID); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✓ Logged out successfully")
		}

		fmt.Println("7. Verifying session after logout...")
		if _, err := a.userManager.ValidateSession(sessionID); err != nil {
			fmt.Println("✓ Session correctly invalidated")
		}
	}

	fmt.Println("8. All registered users:")
	for _, username := range a.userManager.ListAllUsers() {
		fmt.Printf("  - %s\n", username)
	}
}
