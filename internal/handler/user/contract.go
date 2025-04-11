package user

import "github.com/developeerz/restorio-auth/internal/handler/user/dto"

type Service interface {
	SignUp(req *dto.SignUpRequest) (int, error)
	Verify(req *dto.VerificationRequest) (int, error)
	Login(req *dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error)
}
