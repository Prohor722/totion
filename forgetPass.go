package main

import "errors, fmt"

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
	user, exists := s.auth.FindUserByEmail(email)
	if !exists {
		return "Error: email not found"
	}
	// In a real implementation, generate a secure token and send an email
	fmt.Printf("✓ Password reset requested for '%s'. (Token: dummy-token)\n", user.Username)
	return ""
}

func (s *forgetPasswordService) ResetPassword(token, newPassword string) string {
	// In a real implementation, validate the token and reset the password
	if token != "dummy-token" {
		return "Error: invalid token"
	}
	if err := s.auth.ResetPasswordWithToken(token, newPassword); err != nil {
		return "Error: " + err.Error()
	}
	fmt.Println("✓ Password reset successfully")
	return ""
}