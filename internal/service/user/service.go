package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/models"
	"github.com/developeerz/restorio-auth/internal/service/user/mapper"
	"golang.org/x/crypto/bcrypt"
)

const (
	minVal   int = 1000000
	rangeVal int = 8999999
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) SignUp(req *dto.SignUpRequest) (int, error) {
	var err error
	var user *models.User

	user, err = mapper.SignUpToUser(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = service.repo.CreateUser(user)
	if err != nil {
		return http.StatusConflict, err
	}

	userCode := models.UserCode{Telegram: user.Telegram, Code: rand.Intn(rangeVal) + minVal}

	err = service.repo.CreateVerificationCode(&userCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (service *Service) Verify(req *dto.VerificationRequest) (int, error) {
	var err error

	userCode, err := mapper.VerificationToUserCode(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userID, err := service.repo.CheckVerificationCode(userCode)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	err = service.repo.DeleteVerificationCode(userCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAuth := &models.UserAuth{UserID: userID, AuthID: models.USER}

	err = service.repo.SetUserAuth(userAuth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (service *Service) Login(req *dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error) {
	var err error
	var user *models.UserWithAuths

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	user, err = service.repo.FindByTelegramWithAuths(req.Telegram)
	if err != nil {
		return http.StatusNotFound, nil, "", err
	}

	if user.Auths == nil {
		return http.StatusUnauthorized, nil, "", fmt.Errorf("UserID (%d): hasn't auths", user.ID)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return http.StatusUnauthorized, nil, "", err
	}

	jwts, err := jwt.NewJwt(user.ID, mapper.AuthsToStrings(user.Auths))
	if err != nil {
		return http.StatusInternalServerError, nil, "", err
	}

	return http.StatusOK, mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
