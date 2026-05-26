package main

import "errors"

// ForgetPasswordService defines operations for password reset
type ForgetPasswordService interface {
	RequestReset(email string) string
	ResetPassword(token, newPassword string) string
}