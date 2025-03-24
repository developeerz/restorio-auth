package service

import (
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/service/mapper"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (authService *AuthService) Refresh(refreshToken string) (*jwt.Jwt, error) {
	userId, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, err
	}

	userAuths, err := authService.userRepository.GetUserAuths(userId)
	if err != nil {
		return nil, err
	}

	jwt, err := jwt.NewJwt(mapper.UserAuthToIdAndAuth(userAuths))
	if err != nil {
		return nil, err
	}

	return jwt, nil
}
