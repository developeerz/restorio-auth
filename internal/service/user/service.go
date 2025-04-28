package user

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"github.com/developeerz/restorio-auth/internal/service/user/mapper"
	redis "github.com/developeerz/restorio-auth/pkg/repository/redis"
	"golang.org/x/crypto/bcrypt"
)

const (
	minVal   int = 1000000
	rangeVal int = 8999999
)

type Service struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (service *Service) SignUp(ctx context.Context, req *dto.SignUpRequest) (int, error) {
	var err error

	user := mapper.SignUpToUser(req)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = service.cache.PutUser(ctx, req.Telegram, userBytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	verificationCode := genVerificationCode()

	err = service.cache.PutVerificationCode(ctx, req.Telegram, verificationCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (service *Service) Verify(ctx context.Context, req *dto.VerificationRequest) (int, error) {
	var err error
	var userCached redis.User

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	verificationCode, err := service.cache.GetVerificationCode(ctx, req.Telegram)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if verificationCode != req.Code {
		return http.StatusBadRequest,
			fmt.Errorf("wrong verification code: expected: %d, but got: %d", verificationCode, req.Code)
	}

	userByte, err := service.cache.GetUser(ctx, req.Telegram)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(userByte, &userCached)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user, err := mapper.UserToUser(&userCached)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = service.repo.CreateUser(ctx, user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAuth := &models.UserAuth{UserTelegramID: user.TelegramID, AuthID: models.USER}

	err = service.repo.CreateUserAuth(ctx, userAuth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (service *Service) Login(ctx context.Context, req *dto.LoginRequest) (int, *dto.JwtAccessResponse, string, error) {
	var err error
	var user *models.UserWithAuths

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	user, err = service.repo.FindByTelegramWithAuths(ctx, req.Telegram)
	if err != nil {
		return http.StatusNotFound, nil, "", err
	}

	if user.Auths == nil {
		return http.StatusUnauthorized, nil, "", fmt.Errorf("user TelegramID (%d): hasn't auths", user.TelegramID)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return http.StatusUnauthorized, nil, "", err
	}

	jwts, err := jwt.NewJwt(user.TelegramID, mapper.AuthsToStrings(user.Auths))
	if err != nil {
		return http.StatusInternalServerError, nil, "", err
	}

	return http.StatusOK, mapper.JwtToAccess(jwts), jwts.Refresh, nil
}

func genVerificationCode() int {
	return rand.Intn(rangeVal) + minVal
}
