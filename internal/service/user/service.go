package user

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"github.com/developeerz/restorio-auth/internal/service/user/mapper"
	redis_dto "github.com/developeerz/restorio-auth/pkg/redis"
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

func (service *Service) SignUp(req *dto.SignUpRequest) (int, error) {
	var err error

	user := mapper.SignUpToUser(req)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = service.cache.PutUser(req.Telegram, userBytes)
	if err != nil {
		return http.StatusConflict, err
	}

	verificationCode := genVerificationCode()

	err = service.cache.PutVerificationCode(req.Telegram, verificationCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (service *Service) Verify(req *dto.VerificationRequest) (int, error) {
	var err error
	var userCached *redis_dto.User

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	verificationCode, err := service.cache.GetVerificationCode(req.Telegram)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if verificationCode != req.Code {
		return http.StatusBadRequest,
			fmt.Errorf("wrong verification code: expected: %d, but got: %d", verificationCode, req.Code)
	}

	userByte, err := service.cache.GetUser(req.Telegram)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(userByte, userCached)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user, err := mapper.UserToUser(userCached)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = service.repo.CreateUser(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAuth := &models.UserAuth{TelegramID: user.TelegramID, AuthID: models.USER}

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
		return http.StatusUnauthorized, nil, "", fmt.Errorf("UserID (%d): hasn't auths", user.TelegramID)
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
