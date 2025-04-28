package user

import (
	"context"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
)

type Service interface {
	SignUp(ctx context.Context, req *dto.SignUpRequest) (int, error)
	Verify(ctx context.Context, req *dto.VerificationRequest) (int, error)
	Login(ctx context.Context, req *dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error)
}
