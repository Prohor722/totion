package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProcessTerminalInput handles user input from the terminal
func ProcessTerminalInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the User Management System")
	fmt.Println("Available commands: register, login, logout, info, list, delete, changepassword, exit")
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
		command := args[0]
		switch command {
		case "register":
			if len(args) != 4 {
				fmt.Println("Usage: register <username> <email> <password>")
				continue
			}
			fmt.Println(RegisterUser(args[1], args[2], args[3]))
		case "login":
			if len(args) != 3 {
				fmt.Println("Usage: login <username> <password>")
				continue
			}
			sessionID, err := LoginUser(args[1], args[2])
			if err != "" {
				fmt.Println(err)
			} else {
				fmt.Printf("Logged in successfully. Session ID: %s\n", sessionID)
			}
		case "logout":
			if len(args) != 2 {
				fmt.Println("Usage: logout <sessionID>")
				continue
			}

			fmt.Println(LogoutUser(args[1]))
		case "info":
			if len(args) != 2 {
				fmt.Println("Usage: info <sessionID>")
				continue
			}
			user, err := GetUserInfo(args[1])
			if err != "" {
				fmt.Println(err)
			} else {
				fmt.Printf("Username: %s, Email: %s\n", user.Username, user.Email)
			}
		case "list":
			fmt.Println("Registered users:")
			for _, username := range ListAllUsers() {
				fmt.Printf("  - %s\n", username)
			}
		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <username>")
				continue
			}
			fmt.Println(DeleteUser(args[1]))
		case "changepassword":
			if len(args) != 4 {
				fmt.Println("Usage: changepassword <sessionID> <oldPassword> <newPassword>")
				continue
			}
			fmt.Println(ChangePassword(args[1], args[2], args[3]))
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown command. Available commands: register, login, logout, info, list, delete, changepassword, exit")
		}
	}
}