package main

import (
	"os"

	"github.com/Prohor722/totion/auth"
	"github.com/Prohor722/totion/cli"
	"github.com/Prohor722/totion/store"
)

func main() {
	authService := auth.NewAuthService(store.UserStore, auth.NewInMemorySessionRepository())
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cli.ProcessTerminalInputWithAuth(authService)
		return
	}

	app := NewAppWithAuth(authService, store.UserStore)
	app.RunDemo()
}
