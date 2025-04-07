package auth

import (
	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/service/mapper"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) Refresh(refreshToken string) (*dto.JwtAccess, string, error) {
	userID, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, "", err
	}

	userAuths, err := service.repo.GetUserAuths(userID)
	if err != nil {
		return nil, "", err
	}

	jwts, err := jwt.NewJwt(mapper.UserAuthToIDAndAuth(userAuths))
	if err != nil {
		return nil, "", err
	}

	return mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
