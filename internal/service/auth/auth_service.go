package auth

import (
	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/service/mapper"
)

type AuthService struct {
	authRepository AuthRepository
}

func NewAuthService(authRepository AuthRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}

func (authService *AuthService) Refresh(refreshToken string) (*dto.JwtAccess, string, error) {
	userId, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, "", err
	}

	userAuths, err := authService.authRepository.GetUserAuths(userId)
	if err != nil {
		return nil, "", err
	}

	jwts, err := jwt.NewJwt(mapper.UserAuthToIdAndAuth(userAuths))
	if err != nil {
		return nil, "", err
	}

	return mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
