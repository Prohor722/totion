package main

import (
	"fmt"
)

// ForgetPasswordService defines operations for password reset
type ForgetPasswordService interface {
	RequestReset(email string) string
	ResetPassword(token, newPassword string) string
}

type forgetPasswordService struct {
	auth AuthService
}

func NewForgetPasswordService(auth AuthService) ForgetPasswordService {
	return &forgetPasswordService{auth: auth}
}

func (s *forgetPasswordService) RequestReset(email string) string {
	token, err := s.auth.CreatePasswordResetToken(email)
	if err != nil {
		return "Error: " + err.Error()
	}
	fmt.Printf("✓ Password reset requested for '%s'. (Token: %s)\n", email, token)
	return ""
}

func (s *forgetPasswordService) ResetPassword(token, newPassword string) string {
	if err := s.auth.ResetPasswordWithToken(token, newPassword); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Println("✓ Password reset successfully")
	return ""
}
