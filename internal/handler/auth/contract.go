package auth

import (
	"context"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
)

type Service interface {
	Refresh(ctx context.Context, refreshToken string) (*dto.JwtAccessResponse, string, error)
}
