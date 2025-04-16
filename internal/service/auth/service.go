package auth

import (
	user_dto "github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/service/auth/mapper"
	user_mapper "github.com/developeerz/restorio-auth/internal/service/user/mapper"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) Refresh(refreshToken string) (*user_dto.JwtAccessResponse, string, error) {
	telegramID, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, "", err
	}

	userAuths, err := service.repo.GetUserAuths(telegramID)
	if err != nil {
		return nil, "", err
	}

	jwts, err := jwt.NewJwt(mapper.UserAuthToIDAndAuth(userAuths))
	if err != nil {
		return nil, "", err
	}

	return user_mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
