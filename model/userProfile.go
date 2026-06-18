package model

type User struct {
	Username     string
	Email        string
	PasswordHash string
	Profile      *UserProfile
}

type UserProfile struct {
	Username string
	Email    string
	Bio      string
	Website  string
}
