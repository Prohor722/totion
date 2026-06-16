package main

import (
	"os"

	"github.com/Prohor722/totion/auth"
	"github.com/Prohor722/totion/cli"
	"github.com/Prohor722/totion/internal/app"
	"github.com/Prohor722/totion/store"
)

func main() {
	userRepo := store.NewInMemoryUserRepository()
	authService := auth.NewAuthService(userRepo, auth.NewInMemorySessionRepository())
	profileService := auth.NewProfileService(userRepo)

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cli.ProcessTerminalInputWithAuth(authService, profileService)
		return
	}

	application := app.NewAppWithAuth(authService, userRepo)
	application.RunDemo()
}
