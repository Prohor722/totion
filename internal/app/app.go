package app

import (
	"fmt"

	"github.com/Prohor722/totion/auth"
	"github.com/Prohor722/totion/store"
)

type App struct {
	auth           auth.AuthManager
	userManager    UserManager
	profileManager auth.ProfileManager
	resetService   auth.ForgetPasswordService
}

func NewApp(a auth.AuthManager, userManager UserManager, profileManager auth.ProfileManager, resetService auth.ForgetPasswordService) *App {
	return &App{
		auth:           a,
		userManager:    userManager,
		profileManager: profileManager,
		resetService:   resetService,
	}
}

func NewAppWithAuth(authService auth.AuthService, userRepo store.UserRepository) *App {
	return NewApp(
		auth.NewDefaultAuthManager(authService, authService),
		NewDefaultUserManager(authService, authService, authService),
		auth.NewDefaultProfileManager(auth.NewProfileService(userRepo)),
		auth.NewForgetPasswordService(authService),
	)
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

	fmt.Println("9. Viewing and updating profile for john_doe...")
	profile, err := a.profileManager.GetUserProfile("john_doe")
	if err != nil {
		fmt.Println("GetUserProfile error:", err)
	} else {
		fmt.Printf("Current profile: %+v\n", profile)
	}

	if err := a.profileManager.UpdateUserProfile("john_doe", "john.new@example.com", "Updated bio for John"); err != nil {
		fmt.Println("UpdateUserProfile error:", err)
	} else {
		newProfile, _ := a.profileManager.GetUserProfile("john_doe")
		fmt.Printf("Updated profile: %+v\n", newProfile)
	}

	fmt.Println("10. Password reset flow for jane_smith...")
	token, err := a.resetService.RequestReset("jane@example.com")
	if err != nil {
		fmt.Println("RequestReset error:", err)
	} else {
		fmt.Printf("Reset token created for jane_smith: %s\n", token)
		if err := a.resetService.ResetPassword(token, "SafePass1$"); err != nil {
			fmt.Println("ResetPassword error:", err)
		} else {
			fmt.Println("Password reset succeeded for jane_smith")
			if _, err := a.auth.LoginUser("jane_smith", "SafePass1$"); err != nil {
				fmt.Println("Login with new password failed:", err)
			} else {
				fmt.Println("Login with new password succeeded")
			}
		}
	}
}
