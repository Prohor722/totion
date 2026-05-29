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
	if msg := a.auth.RegisterUser("john_doe", "john@example.com", "password123"); msg != "" {
		fmt.Println(msg)
	}
	if msg := a.auth.RegisterUser("jane_smith", "jane@example.com", "securePass456"); msg != "" {
		fmt.Println(msg)
	}

	fmt.Println(a.auth.RegisterUser("ab", "invalid@email.com", "short"))

	fmt.Println("2. Login attempt...")
	sessionID, errMsg := a.auth.LoginUser("john_doe", "password123")
	if errMsg != "" {
		fmt.Println("Login error:", errMsg)
	} else {
		fmt.Printf("Session ID: %s\n", sessionID)

		fmt.Println("3. Validating session...")
		if valid, username := a.userManager.ValidateSession(sessionID); valid {
			fmt.Printf("✓ Session is valid for user: %s\n", username)
		}

		fmt.Println("4. Retrieving user info...")
		user, msg := a.userManager.GetUserInfo(sessionID)
		if msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
		}

		fmt.Println("5. Changing password...")
		if msg := a.userManager.ChangePassword(sessionID, "password123", "newPassword789"); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Println("✓ Password changed successfully")
		}

		fmt.Println("6. Logging out...")
		if msg := a.auth.LogoutUser(sessionID); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Println("✓ Logged out successfully")
		}

		fmt.Println("7. Verifying session after logout...")
		if valid, _ := a.userManager.ValidateSession(sessionID); !valid {
			fmt.Println("✓ Session correctly invalidated")
		}
	}

	fmt.Println("8. All registered users:")
	for _, username := range a.userManager.ListAllUsers() {
		fmt.Printf("  - %s\n", username)
	}
}
