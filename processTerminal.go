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