package main

import (
	"sort"
	"sync"
)

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
	mutex sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserRepository creates a new in-memory repo
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*User)}
}

func (r *InMemoryUserRepository) Get(username string) (*User, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	u, ok := r.users[username]
	return u, ok
}

func (r *InMemoryUserRepository) Add(user *User) error {
	if user == nil {
		return ErrUserIsNil
	}

    r.mutex.Lock()
    defer r.mutex.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return ErrUserExists
	}
    r.users[user.Username] = user
    return nil
}

func (r *InMemoryUserRepository) Exists(username string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, ok := r.users[username]
	return ok
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*User, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, true
		}
	}
	return nil, false
}

func (r *InMemoryUserRepository) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var out []string
	for k := range r.users {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (r *InMemoryUserRepository) Delete(username string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.users[username]; !ok {
		return ErrUserNotFound
	}
	delete(r.users, username)
	return nil
}

// Package-level default repository to gradually replace globals
var UserStore UserRepository = NewInMemoryUserRepository()
