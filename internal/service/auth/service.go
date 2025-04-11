package auth

import (
	user_dto "github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	user_mapper "github.com/developeerz/restorio-auth/internal/service/user/mapper"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) Refresh(refreshToken string) (*user_dto.JwtAccessResponse, string, error) {
	userID, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, "", err
	}

	userAuths, err := service.repo.GetUserAuths(userID)
	if err != nil {
		return nil, "", err
	}

	jwts, err := jwt.NewJwt(user_mapper.UserAuthToIDAndAuth(userAuths))
	if err != nil {
		return nil, "", err
	}

	return user_mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
