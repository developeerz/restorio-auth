package auth

import "github.com/developeerz/restorio-auth/internal/handler/user/dto"

type Service interface {
	Refresh(refreshToken string) (*dto.JwtAccessResponse, string, error)
}
