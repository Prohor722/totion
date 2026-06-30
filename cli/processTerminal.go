package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Prohor722/totion/auth"
	"github.com/Prohor722/totion/model"
	"github.com/Prohor722/totion/store"
)

const genericPasswordResetResponse = "If an account with that email exists, password reset instructions have been sent."

// Command represents a terminal action
type Command interface {
	Execute(args []string) (string, error)
}

// UserService defines user-related operations used by commands
type UserService interface {
	Register(username, email, password string) error
	ListAll() []string
	Delete(username string) error
	ChangePassword(sessionID, oldPassword, newPassword string) error
	GetInfo(sessionID string) (*model.User, error)
}

// SessionService defines session-related operations used by commands
type SessionService interface {
	Login(username, password string) (string, error)
	Logout(sessionID string) error
}

type PasswordResetService interface {
	RequestReset(email string) (string, error)
	ResetPassword(token, newPassword string) error
}

type ProfileService interface {
	UpdateProfile(username string, update auth.ProfileUpdate) error
}

type authUserService struct {
	account  auth.RegistrationService
	password auth.PasswordService
	session  auth.SessionValidationService
}

type authResetService struct {
	reset auth.ForgetPasswordService
}

type authProfileService struct {
	profile auth.ProfileService
}

func NewTerminalPasswordResetService(reset auth.ForgetPasswordService) PasswordResetService {
	return &authResetService{reset: reset}
}

func NewTerminalProfileService(profile auth.ProfileService) ProfileService {
	return &authProfileService{profile: profile}
}

func (d *authResetService) RequestReset(email string) (string, error) {
	return d.reset.RequestReset(email)
}

func (d *authResetService) ResetPassword(token, newPassword string) error {
	return d.reset.ResetPassword(token, newPassword)
}

func (d *authProfileService) UpdateProfile(username string, update auth.ProfileUpdate) error {
	return d.profile.UpdateProfile(username, update)
}

func NewTerminalUserService(account auth.RegistrationService, password auth.PasswordService, session auth.SessionValidationService) UserService {
	return &authUserService{account: account, password: password, session: session}
}

func (d *authUserService) Register(u, e, p string) error {
	return d.account.Register(u, e, p)
}

func (d *authUserService) ListAll() []string {
	return d.account.ListUsernames()
}

func (d *authUserService) Delete(u string) error {
	return d.account.DeleteUser(u)
}

func (d *authUserService) ChangePassword(s, o, n string) error {
	return d.password.ChangePassword(s, o, n)
}

func (d *authUserService) GetInfo(sessionID string) (*model.User, error) {
	return d.session.GetUserInfo(sessionID)
}

type authSessionService struct {
	credentials auth.CredentialService
}

func NewTerminalSessionService(credentials auth.CredentialService) SessionService {
	return &authSessionService{credentials: credentials}
}

func (d *authSessionService) Login(u, p string) (string, error) {
	return d.credentials.Login(u, p)
}

func (d *authSessionService) Logout(s string) error {
	return d.credentials.Logout(s)
}

// Concrete command implementations
type registerCommand struct{ users UserService }

func (c *registerCommand) Execute(args []string) (string, error) {
	if len(args) != 4 {
		return "", errors.New("Usage: register <username> <email> <password>")
	}
	if err := c.users.Register(args[1], args[2], args[3]); err != nil {
		return "", err
	}
	return "Registered successfully", nil
}

type loginCommand struct{ sessions SessionService }

func (c *loginCommand) Execute(args []string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("Usage: login <username> <password>")
	}
	return c.sessions.Login(args[1], args[2])
}

type logoutCommand struct{ sessions SessionService }

func (c *logoutCommand) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Usage: logout <sessionID>")
	}
	if err := c.sessions.Logout(args[1]); err != nil {
		return "", err
	}
	return "Logged out successfully", nil
}

type infoCommand struct{ users UserService }

func (c *infoCommand) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Usage: info <sessionID>")
	}
	user, err := c.users.GetInfo(args[1])
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	return fmt.Sprintf("Username: %s, Email: %s", user.Username, user.Email), nil
}

type listCommand struct{ users UserService }

func (c *listCommand) Execute(args []string) (string, error) {
	users := c.users.ListAll()
	if len(users) == 0 {
		return "No registered users.", nil
	}
	var b strings.Builder
	b.WriteString("Registered users:\n")
	for _, u := range users {
		b.WriteString("  - ")
		b.WriteString(u)
		b.WriteByte('\n')
	}
	return b.String(), nil
}

type deleteCommand struct{ users UserService }

func (c *deleteCommand) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Usage: delete <username>")
	}
	if err := c.users.Delete(args[1]); err != nil {
		return "", err
	}
	return "User deleted successfully", nil
}

type changePasswordCommand struct{ users UserService }

func (c *changePasswordCommand) Execute(args []string) (string, error) {
	if len(args) != 4 {
		return "", errors.New("Usage: changepassword <sessionID> <oldPassword> <newPassword>")
	}
	if err := c.users.ChangePassword(args[1], args[2], args[3]); err != nil {
		return "", err
	}
	return "Password changed successfully", nil
}

type requestResetCommand struct{ reset PasswordResetService }

func (c *requestResetCommand) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Usage: requestreset <email>")
	}

	_, err := c.reset.RequestReset(args[1])
	if err != nil {
		if errors.Is(err, auth.ErrEmailNotFound) || errors.Is(err, auth.ErrTooManyResetRequests) {
			return genericPasswordResetResponse, nil
		}
		return "", err
	}

	return genericPasswordResetResponse, nil
}

type resetPasswordCommand struct{ reset PasswordResetService }

func (c *resetPasswordCommand) Execute(args []string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("Usage: resetpassword <token> <newPassword>")
	}
	if err := c.reset.ResetPassword(args[1], args[2]); err != nil {
		return "", err
	}
	return "Password reset successfully", nil
}

type viewProfileCommand struct{ users UserService }

func (c *viewProfileCommand) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Usage: viewprofile <sessionID>")
	}
	user, err := c.users.GetInfo(args[1])
	if err != nil {
		return "", err
	}
	if user.Profile == nil {
		return "Profile not available", nil
	}
	result := fmt.Sprintf("Profile for %s:\n  Email: %s\n  Bio: %s", user.Username, user.Profile.Email, user.Profile.Bio)
	if user.Profile.Website != "" {
		result += fmt.Sprintf("\n  Website: %s", user.Profile.Website)
	}
	return result, nil
}

type updateProfileCommand struct {
	users    UserService
	profiles ProfileService
}


// ProcessTerminalInput handles user input from the terminal using a small, testable processor
func ProcessTerminalInputWithAuth(authService auth.AuthService, profileService auth.ProfileService) {
	ProcessTerminalInputWithServices(
		NewTerminalUserService(authService, authService, authService),
		NewTerminalSessionService(authService),
		NewTerminalPasswordResetService(auth.NewForgetPasswordService(authService)),
		NewTerminalProfileService(profileService),
	)
}

func ProcessTerminalInputWithServices(users UserService, sessions SessionService, reset PasswordResetService, profiles ProfileService) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the User Management System")
	fmt.Println("Available commands: register, login, logout, info, list, delete, changepassword, requestreset, resetpassword, viewprofile, updateprofile, exit")

	commands := map[string]Command{
		"register":       &registerCommand{users: users},
		"login":          &loginCommand{sessions: sessions},
		"logout":         &logoutCommand{sessions: sessions},
		"info":           &infoCommand{users: users},
		"list":           &listCommand{users: users},
		"delete":         &deleteCommand{users: users},
		"changepassword": &changePasswordCommand{users: users},
		"requestreset":   &requestResetCommand{reset: reset},
		"resetpassword":  &resetPasswordCommand{reset: reset},
		"viewprofile":    &viewProfileCommand{users: users},
		"updateprofile":  &updateProfileCommand{users: users, profiles: profiles},
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		args := strings.Fields(input)
		cmdName := args[0]
		if cmdName == "exit" {
			fmt.Println("Exiting...")
			return
		}
		cmd, ok := commands[cmdName]
		if !ok {
			fmt.Println("Unknown command. Available commands: register, login, logout, info, list, delete, changepassword, requestreset, resetpassword, viewprofile, updateprofile, exit")
			continue
		}
		result, err := cmd.Execute(args)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if result != "" {
			fmt.Println(result)
		}
	}
}

// ProcessTerminalInput starts the interactive CLI with the default auth and profile services.
func ProcessTerminalInput() {
	ProcessTerminalInputWithAuth(auth.DefaultAuth, auth.NewProfileService(store.UserStore))
}
