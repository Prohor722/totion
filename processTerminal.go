package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Command represents a terminal action
type Command interface {
	Execute(args []string) string
}

// UserService defines user-related operations used by commands
type UserService interface {
	Register(username, email, password string) error
	ListAll() []string
	Delete(username string) error
	ChangePassword(sessionID, oldPassword, newPassword string) error
	GetInfo(sessionID string) (*User, error)
}

// SessionService defines session-related operations used by commands
type SessionService interface {
	Login(username, password string) (string, error)
	Logout(sessionID string) error
}

type authUserService struct {
	auth AuthService
}

func (d *authUserService) Register(u, e, p string) error {
	return d.auth.Register(u, e, p)
}

func (d *authUserService) ListAll() []string {
	return d.auth.ListUsernames()
}

func (d *authUserService) Delete(u string) error {
	return d.auth.DeleteUser(u)
}

func (d *authUserService) ChangePassword(s, o, n string) error {
	return d.auth.ChangePassword(s, o, n)
}

func (d *authUserService) GetInfo(sessionID string) (*User, error) {
	return d.auth.GetUserInfo(sessionID)
}

type authSessionService struct {
	auth AuthService
}

func (d *authSessionService) Login(u, p string) (string, error) {
	return d.auth.Login(u, p)
}

func (d *authSessionService) Logout(s string) error {
	return d.auth.Logout(s)
}

// Concrete command implementations
type registerCommand struct{ users UserService }

func (c *registerCommand) Execute(args []string) string {
	if len(args) != 4 {
		return "Usage: register <username> <email> <password>"
	}
	if err := c.users.Register(args[1], args[2], args[3]); err != nil {
		return "Error: " + err.Error()
	}
	return "Registered successfully"
}

type loginCommand struct{ sessions SessionService }

func (c *loginCommand) Execute(args []string) string {
	if len(args) != 3 {
		return "Usage: login <username> <password>"
	}
	sessionID, err := c.sessions.Login(args[1], args[2])
	if err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Logged in successfully. Session ID: %s", sessionID)
}

type logoutCommand struct{ sessions SessionService }

func (c *logoutCommand) Execute(args []string) string {
	if len(args) != 2 {
		return "Usage: logout <sessionID>"
	}
	if err := c.sessions.Logout(args[1]); err != nil {
		return "Error: " + err.Error()
	}
	return "Logged out successfully"
}

type infoCommand struct{ users UserService }

func (c *infoCommand) Execute(args []string) string {
	if len(args) != 2 {
		return "Usage: info <sessionID>"
	}
	user, err := c.users.GetInfo(args[1])
	if err != nil {
		return "Error: " + err.Error()
	}
	if user == nil {
		return "Error: User not found"
	}
	return fmt.Sprintf("Username: %s, Email: %s", user.Username, user.Email)
}

type listCommand struct{ users UserService }

func (c *listCommand) Execute(args []string) string {
	users := c.users.ListAll()
	if len(users) == 0 {
		return "No registered users."
	}
	var b strings.Builder
	b.WriteString("Registered users:\n")
	for _, u := range users {
		b.WriteString("  - ")
		b.WriteString(u)
		b.WriteByte('\n')
	}
	return b.String()
}

type deleteCommand struct{ users UserService }

func (c *deleteCommand) Execute(args []string) string {
	if len(args) != 2 {
		return "Usage: delete <username>"
	}
	if err := c.users.Delete(args[1]); err != nil {
		return "Error: " + err.Error()
	}
	return "User deleted successfully"
}

type changePasswordCommand struct{ users UserService }

func (c *changePasswordCommand) Execute(args []string) string {
	if len(args) != 4 {
		return "Usage: changepassword <sessionID> <oldPassword> <newPassword>"
	}
	if err := c.users.ChangePassword(args[1], args[2], args[3]); err != nil {
		return "Error: " + err.Error()
	}
	return "Password changed successfully"
}

// ProcessTerminalInput handles user input from the terminal using a small, testable processor
func ProcessTerminalInputWithAuth(auth AuthService) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the User Management System")
	fmt.Println("Available commands: register, login, logout, info, list, delete, changepassword, exit")

	// build services and command registry (dependencies injected)
	us := &authUserService{auth: auth}
	ss := &authSessionService{auth: auth}

	commands := map[string]Command{
		"register":       &registerCommand{users: us},
		"login":          &loginCommand{sessions: ss},
		"logout":         &logoutCommand{sessions: ss},
		"info":           &infoCommand{users: us},
		"list":           &listCommand{users: us},
		"delete":         &deleteCommand{users: us},
		"changepassword": &changePasswordCommand{users: us},
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
		args := strings.Split(input, " ")
		cmdName := args[0]
		if cmdName == "exit" {
			fmt.Println("Exiting...")
			return
		}
		cmd, ok := commands[cmdName]
		if !ok {
			fmt.Println("Unknown command. Available commands: register, login, logout, info, list, delete, changepassword, exit")
			continue
		}
		result := cmd.Execute(args)
		if result != "" {
			fmt.Println(result)
		}
	}
}

// ProcessTerminalInput starts the interactive CLI with the default auth service.
func ProcessTerminalInput() {
	ProcessTerminalInputWithAuth(DefaultAuth)
}
