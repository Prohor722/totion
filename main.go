package main

func main() {
	authService := NewAuthService(UserStore, NewInMemorySessionRepository())
	app := NewAppWithAuth(authService, UserStore)
	app.RunDemo()
}
