package main

import "errors"

// UserRepository defines operations for user persistence
type UserRepository interface {
	Get(username string) (*User, bool)
	Add(user *User) error
	Exists(username string) bool
	FindByEmail(email string) (*User, bool)
	List() []string
	Delete(username string) error
}

// InMemoryUserRepository is a simple in-memory store for users
type InMemoryUserRepository struct {
	users map[string]*User
}

// NewInMemoryUserRepository creates a new in-memory repo
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*User)}
}

func (r *InMemoryUserRepository) Get(username string) (*User, bool) {
	u, ok := r.users[username]
	return u, ok
}

func (r *InMemoryUserRepository) Add(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if _, exists := r.users[user.Username]; exists {
		return errors.New("username already exists")
	}
	r.users[user.Username] = user
	return nil
}

func (r *InMemoryUserRepository) Exists(username string) bool {
	_, ok := r.users[username]
	return ok
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*User, bool) {
	for _, u := range r.users {
		if u.Email == email {
			return u, true
		}
	}
	return nil, false
}

func (r *InMemoryUserRepository) List() []string {
	var out []string
	for k := range r.users {
		out = append(out, k)
	}
	return out
}

func (r *InMemoryUserRepository) Delete(username string) error {
	if _, ok := r.users[username]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, username)
	return nil
}

// Package-level default repository to gradually replace globals
var UserStore UserRepository = NewInMemoryUserRepository()
