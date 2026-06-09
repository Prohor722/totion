package main

// ForgetPasswordService defines operations for password reset
type ForgetPasswordService interface {
	RequestReset(email string) (string, error)
	ResetPassword(token, newPassword string) error
}

type forgetPasswordService struct {
	password PasswordService
}

func NewForgetPasswordService(auth PasswordService) ForgetPasswordService {
	return &forgetPasswordService{password: auth}
}

func (s *forgetPasswordService) RequestReset(email string) (string, error) {
	return s.password.CreatePasswordResetToken(email)
}

func (s *forgetPasswordService) ResetPassword(token, newPassword string) error {
	return s.password.ResetPasswordWithToken(token, newPassword)
}
