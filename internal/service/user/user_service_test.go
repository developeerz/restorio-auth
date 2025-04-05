package user_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/models"
	"github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/developeerz/restorio-auth/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpSuccess(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	mockUserRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
	mockUserRepo.On("CreateVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(nil)

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	status, err := service.SignUp(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	mockUserRepo.AssertExpectations(t)
}

func TestSignUpDublicatedUser(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	mockUserRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(errors.New("Dublicate user"))

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	status, err := service.SignUp(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusConflict, status)

	mockUserRepo.AssertExpectations(t)
}

func TestSignUpDublicatedVerificationCode(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	mockUserRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
	mockUserRepo.On("CreateVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(errors.New("Dublicate verification code"))

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	status, err := service.SignUp(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)

	mockUserRepo.AssertExpectations(t)
}

func TestVerifySuccess(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	var userId int64 = 1

	mockUserRepo.On("CheckVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(userId, nil)
	mockUserRepo.On("DeleteVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(nil)
	mockUserRepo.On("VerifyUser", userId).Return(nil)
	mockUserRepo.On("SetUserAuth", mock.AnythingOfType("*models.UserAuth")).Return(nil)

	req := &dto.VerificationRequest{Code: 111000, Telegram: "@telegram"}
	code, err := service.Verify(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	mockUserRepo.AssertExpectations(t)
}

func TestVerifyWrongVerificationCode(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	mockUserRepo.On("CheckVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(int64(0), errors.New("Wrong verification code"))

	req := &dto.VerificationRequest{Code: 111000, Telegram: "@telegram"}
	code, err := service.Verify(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)

	mockUserRepo.AssertExpectations(t)
}

func TestVerifyDeleteVerificationCodeError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	var userId int64 = 1

	mockUserRepo.On("CheckVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(userId, nil)
	mockUserRepo.On("DeleteVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(errors.New("Cannot delete verification code"))

	req := &dto.VerificationRequest{Code: 111000, Telegram: "@telegram"}
	code, err := service.Verify(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserRepo.AssertExpectations(t)
}

func TestVerifyVerifyUserError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	var userId int64 = 1

	mockUserRepo.On("CheckVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(userId, nil)
	mockUserRepo.On("DeleteVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(nil)
	mockUserRepo.On("VerifyUser", userId).Return(errors.New("Verification error"))

	req := &dto.VerificationRequest{Code: 111000, Telegram: "@telegram"}
	code, err := service.Verify(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserRepo.AssertExpectations(t)
}

func TestVerifySetUserAuthError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	var userId int64 = 1

	mockUserRepo.On("CheckVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(userId, nil)
	mockUserRepo.On("DeleteVerificationCode", mock.AnythingOfType("*models.UserCode")).Return(nil)
	mockUserRepo.On("VerifyUser", userId).Return(nil)
	mockUserRepo.On("SetUserAuth", mock.AnythingOfType("*models.UserAuth")).Return(errors.New("Cannot create user auth"))

	req := &dto.VerificationRequest{Code: 111000, Telegram: "@telegram"}
	code, err := service.Verify(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	telegram := "telegram"
	decryptedPassword := "password"
	bpassword, _ := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	password := string(bpassword)
	user := &models.User{ID: 1, Telegram: telegram, Password: password, Verified: true}
	userAuths := []models.UserAuth{{UserId: user.ID, AuthId: models.USER}}

	mockUserRepo.On("FindByTelegram", telegram).Return(user, nil)
	mockUserRepo.On("GetUserAuths", user.ID).Return(userAuths, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginFindByTelegramError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	telegram := "telegram"
	decryptedPassword := "password"

	mockUserRepo.On("FindByTelegram", telegram).Return(nil, errors.New("Cannot find user"))

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginNotVerified(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	telegram := "telegram"
	decryptedPassword := "password"
	bpassword, _ := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	password := string(bpassword)
	user := &models.User{ID: 1, Telegram: telegram, Password: password, Verified: false}

	mockUserRepo.On("FindByTelegram", telegram).Return(user, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginWrongPassword(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	telegram := "telegram"
	decryptedPassword := "password"
	bpassword, _ := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	password := string(bpassword)
	user := &models.User{ID: 1, Telegram: telegram, Password: password, Verified: true}

	mockUserRepo.On("FindByTelegram", telegram).Return(user, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: "abracadabra"}
	code, _, _, err := service.Login(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginUserAuthError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := user.NewUserService(mockUserRepo)

	telegram := "telegram"
	decryptedPassword := "password"
	bpassword, _ := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	password := string(bpassword)
	user := &models.User{ID: 1, Telegram: telegram, Password: password, Verified: true}

	mockUserRepo.On("FindByTelegram", telegram).Return(user, nil)
	mockUserRepo.On("GetUserAuths", user.ID).Return(nil, errors.New("Get auths error"))

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserRepo.AssertExpectations(t)
}
