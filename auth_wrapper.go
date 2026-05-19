package main

// DefaultAuth is a package-level default auth service wired with in-memory stores
var DefaultAuth AuthService = NewAuthService(UserStore, NewInMemorySessionRepository())

// LoginUser is a convenience wrapper used by legacy code
func LoginUser(username, password string) (string, string) {
	sid, err := DefaultAuth.Login(username, password)
	if err != nil {
		return "", err.Error()
	}
	return sid, ""
}

// LogoutUser is a convenience wrapper used by legacy code
func LogoutUser(sessionID string) string {
	if err := DefaultAuth.Logout(sessionID); err != nil {
		return err.Error()
	}
	return ""
}
