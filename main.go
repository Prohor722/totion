package main

import "os"

func main() {
	authService := NewAuthService(UserStore, NewInMemorySessionRepository())
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		ProcessTerminalInputWithAuth(authService)
		return
	}

	app := NewAppWithAuth(authService, UserStore)
	app.RunDemo()
}
