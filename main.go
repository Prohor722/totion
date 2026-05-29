package main

func main() {
	authService := NewAuthService(UserStore, NewInMemorySessionRepository())
	app := NewApp(authService, UserStore)
	app.RunDemo()
}
