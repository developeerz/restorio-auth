package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/models"
	"github.com/developeerz/restorio-auth/internal/service/mapper"
	"golang.org/x/crypto/bcrypt"
)

const (
	minVal   int = 1000000
	rangeVal int = 8999999
)

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) *Service {
	return &Service{userRepository: userRepository}
}

func (userService *Service) SignUp(req *dto.SignUpRequest) (int, error) {
	var err error
	var user *models.User

	user, err = mapper.SignUpToUser(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = userService.userRepository.CreateUser(user)
	if err != nil {
		return http.StatusConflict, err
	}

	userCode := models.UserCode{Telegram: user.Telegram, Code: rand.Intn(rangeVal) + minVal}

	err = userService.userRepository.CreateVerificationCode(&userCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (userService *Service) Verify(req *dto.VerificationRequest) (int, error) {
	var err error

	userCode, err := mapper.VerificationToUserCode(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userID, err := userService.userRepository.CheckVerificationCode(userCode)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	err = userService.userRepository.DeleteVerificationCode(userCode)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = userService.userRepository.VerifyUser(userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAuth := &models.UserAuth{UserID: userID, AuthID: models.USER}

	err = userService.userRepository.SetUserAuth(userAuth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (userService *Service) Login(req *dto.LoginRequest) (int, *dto.JwtAccess, string, error) {
	var err error
	var user *models.User

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	user, err = userService.userRepository.FindByTelegram(req.Telegram)
	if err != nil {
		return http.StatusNotFound, nil, "", err
	}

	if !user.Verified {
		return http.StatusUnauthorized, nil, "", fmt.Errorf("GetUserAuths(%d): not verified", user.ID)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return http.StatusUnauthorized, nil, "", err
	}

	userAuths, err := userService.userRepository.GetUserAuths(user.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf("GetUserAuths(%d): %w", user.ID, err)
	}

	jwts, err := jwt.NewJwt(mapper.UserAuthToIDAndAuth(userAuths))
	if err != nil {
		return http.StatusInternalServerError, nil, "", err
	}

	return http.StatusOK, mapper.JwtToAccess(jwts), jwts.Refresh, nil
}
