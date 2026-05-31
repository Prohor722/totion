package main


// ForgetPasswordService defines operations for password reset
type ForgetPasswordService interface {
	RequestReset(email string) (string, error)
	ResetPassword(token, newPassword string) error
}

type forgetPasswordService struct {
	auth AuthService
}

func NewForgetPasswordService(auth AuthService) ForgetPasswordService {
	return &forgetPasswordService{auth: auth}
}

func (s *forgetPasswordService) RequestReset(email string) (string, error) {
	return s.auth.CreatePasswordResetToken(email)
}

func (s *forgetPasswordService) ResetPassword(token, newPassword string) error {
	return s.auth.ResetPasswordWithToken(token, newPassword)
}
