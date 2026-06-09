package main

import "errors"

var (
	ErrUserExists   = errors.New("username already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUserIsNil    = errors.New("user is nil")
	ErrInvalidInput = errors.New("invalid input")
)
