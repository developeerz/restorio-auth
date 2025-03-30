package service

import (
	"math/rand"
	"net/http"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/models"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/service/mapper"
	"golang.org/x/crypto/bcrypt"
)

const (
	minVal   int = 1000000
	rangeVal int = 8999999
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) SignUp(req *dto.SignUpRequest) (int, *dto.Error, error) {
	var err error
	var user *models.User

	user, err = mapper.SignUpToUser(req)
	if err != nil {
		return http.StatusInternalServerError, &dto.Error{Message: "Internal error"}, err
	}

	err = userService.userRepository.CreateUser(user)
	if err != nil {
		return http.StatusConflict, &dto.Error{Message: "User already signed up"}, err
	}

	userCode := models.UserCode{Telegram: user.Telegram, Code: rand.Intn(rangeVal) + minVal}
	err = userService.userRepository.CreateVerificationCode(&userCode)
	if err != nil {
		return http.StatusInternalServerError, &dto.Error{Message: "Internal error"}, err
	}

	return http.StatusOK, nil, nil
}

func (userService *UserService) Verify(req *dto.VerificationRequest) (int, *dto.Error, error) {
	var err error

	userId, err := userService.userRepository.CheckVerificationCode(req.Telegram, req.Code)
	if err != nil {
		return http.StatusUnauthorized, &dto.Error{Message: "Wrong code or telegram"}, err
	}

	userCode := &models.UserCode{Telegram: req.Telegram, Code: req.Code}
	err = userService.userRepository.DeleteVerificationCode(userCode)
	if err != nil {
		return http.StatusInternalServerError, &dto.Error{Message: "Internal error"}, err
	}

	err = userService.userRepository.VerifyUser(userId)
	if err != nil {
		return http.StatusInternalServerError, &dto.Error{Message: "Internal error"}, err
	}

	userAuth := &models.UserAuth{UserId: userId, AuthId: models.USER}
	err = userService.userRepository.SetUserAuth(userAuth)
	if err != nil {
		return http.StatusInternalServerError, &dto.Error{Message: "Internal error"}, err
	}

	return http.StatusOK, nil, nil
}

func (userService *UserService) SignIn(req *dto.SignInRequest) (int, *jwt.Jwt, *dto.Error, error) {
	var err error
	var user *models.User

	user, err = userService.userRepository.FindByTelegram(req.Telegram)
	if err != nil {
		return http.StatusNotFound, nil, &dto.Error{Message: "User not Found"}, err
	}

	if !user.Verified {
		return http.StatusUnauthorized, nil, &dto.Error{Message: "User not verified"}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return http.StatusUnauthorized, nil, &dto.Error{Message: "Wrong Password"}, err
	}

	userAuths, err := userService.userRepository.GetUserAuths(user.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, &dto.Error{Message: "Internal error"}, err
	}

	jwt, err := jwt.NewJwt(mapper.UserAuthToIdAndAuth(userAuths))
	if err != nil {
		return http.StatusInternalServerError, nil, &dto.Error{Message: "Internal error"}, err
	}

	return http.StatusOK, jwt, nil, nil
}
