package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/developeerz/restorio-auth/internal/service/user/mapper"
	redis "github.com/developeerz/restorio-auth/pkg/repository/redis"
	mocks_user "github.com/developeerz/restorio-auth/test/mocks/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const (
	telegram          = "telegram"
	decryptedPassword = "password"
)

func TestSignUpSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	user := mapper.SignUpToUser(req)

	userBytes, err := json.Marshal(user)
	assert.NoError(t, err)

	mockUserCache.On("PutUser", ctx, telegram, userBytes).Return(nil)
	mockUserCache.On("PutVerificationCode", ctx, telegram, mock.AnythingOfType("int")).Return(nil)

	status, err := service.SignUp(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	mockUserCache.AssertExpectations(t)
}

func TestSignUpPutUserError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	user := mapper.SignUpToUser(req)

	userBytes, err := json.Marshal(user)
	assert.NoError(t, err)

	mockUserCache.On("PutUser", ctx, telegram, userBytes).Return(errors.New(""))

	status, err := service.SignUp(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)

	mockUserCache.AssertExpectations(t)
}

func TestSignUpPutVerificationCodeError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	req := &dto.SignUpRequest{
		Firstname: "firstname",
		Lastname:  "lastname",
		Telegram:  "@telegram",
		Password:  "password",
	}
	user := mapper.SignUpToUser(req)

	userBytes, err := json.Marshal(user)
	assert.NoError(t, err)

	mockUserCache.On("PutUser", ctx, telegram, userBytes).Return(nil)
	mockUserCache.On("PutVerificationCode", ctx, telegram, mock.AnythingOfType("int")).Return(errors.New(""))

	status, err := service.SignUp(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)

	mockUserCache.AssertExpectations(t)
}

func TestVerifySuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	var userTelegramID int64 = 111111
	verificationCode := 111111

	userCached := &redis.User{
		TelegramID: userTelegramID,
		Firstname:  "firstname",
		Lastname:   "lastname",
		Telegram:   telegram,
		Password:   decryptedPassword,
	}

	userBytes, err := json.Marshal(userCached)
	assert.NoError(t, err)

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(verificationCode, nil)
	mockUserCache.On("GetUser", ctx, telegram).Return(userBytes, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	mockUserRepo.On("CreateUserAuth", ctx, mock.AnythingOfType("*models.UserAuth")).Return(nil)

	req := &dto.VerificationRequest{Code: verificationCode, Telegram: "@telegram"}
	code, err := service.Verify(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	mockUserCache.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestVerifyWrongTelegram(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	verificationCode := 111111

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(-1, errors.New(""))

	req := &dto.VerificationRequest{Code: verificationCode, Telegram: "@telegram"}

	code, err := service.Verify(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, code)

	mockUserCache.AssertExpectations(t)
}

func TestVerifyWrongVerificationCode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	expectedVerificationCode := 111111
	actualVerificationCode := 999999

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(actualVerificationCode, nil)

	req := &dto.VerificationRequest{Code: expectedVerificationCode, Telegram: "@telegram"}

	code, err := service.Verify(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, code)

	mockUserCache.AssertExpectations(t)
}

func TestVerifyGetUserError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	verificationCode := 111111

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(verificationCode, nil)
	mockUserCache.On("GetUser", ctx, telegram).Return(nil, errors.New(""))

	req := &dto.VerificationRequest{Code: verificationCode, Telegram: "@telegram"}

	code, err := service.Verify(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserCache.AssertExpectations(t)
}

func TestVerifyCreateUserError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	var userTelegramID int64 = 111111
	verificationCode := 111111

	userCached := &redis.User{
		TelegramID: userTelegramID,
		Firstname:  "firstname",
		Lastname:   "lastname",
		Telegram:   telegram,
		Password:   decryptedPassword,
	}

	userBytes, err := json.Marshal(userCached)
	assert.NoError(t, err)

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(verificationCode, nil)
	mockUserCache.On("GetUser", ctx, telegram).Return(userBytes, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(errors.New(""))

	req := &dto.VerificationRequest{Code: verificationCode, Telegram: "@telegram"}

	code, err := service.Verify(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserCache.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestVerifyCreateUserAuthError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	var userTelegramID int64 = 111111
	verificationCode := 111111

	userCached := &redis.User{
		TelegramID: userTelegramID,
		Firstname:  "firstname",
		Lastname:   "lastname",
		Telegram:   telegram,
		Password:   decryptedPassword,
	}

	userBytes, err := json.Marshal(userCached)
	assert.NoError(t, err)

	mockUserCache.On("GetVerificationCode", ctx, telegram).Return(verificationCode, nil)
	mockUserCache.On("GetUser", ctx, telegram).Return(userBytes, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	mockUserRepo.On("CreateUserAuth", ctx, mock.AnythingOfType("*models.UserAuth")).Return(errors.New(""))

	req := &dto.VerificationRequest{Code: verificationCode, Telegram: "@telegram"}

	code, err := service.Verify(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, code)

	mockUserCache.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	bpassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	password := string(bpassword)
	auths := []models.Authority{{ID: "USER", Description: "descr"}}
	user := &models.UserWithAuths{TelegramID: 111111, Telegram: telegram, Password: password, Auths: auths}

	mockUserRepo.On("FindByTelegramWithAuths", ctx, telegram).Return(user, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginFindByTelegramError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	mockUserRepo.On("FindByTelegramWithAuths", ctx, telegram).Return(nil, errors.New("Cannot find user"))

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginNotVerified(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	bpassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	password := string(bpassword)
	user := &models.UserWithAuths{TelegramID: 111111, Telegram: telegram, Password: password, Auths: nil}

	mockUserRepo.On("FindByTelegramWithAuths", ctx, telegram).Return(user, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: decryptedPassword}
	code, _, _, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)

	mockUserRepo.AssertExpectations(t)
}

func TestLoginWrongPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockUserRepo := mocks_user.NewRepository(t)
	mockUserCache := mocks_user.NewCache(t)
	service := user.NewService(mockUserRepo, mockUserCache)

	bpassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	password := string(bpassword)
	auths := []models.Authority{{ID: "USER", Description: "descr"}}
	user := &models.UserWithAuths{TelegramID: 111111, Telegram: telegram, Password: password, Auths: auths}

	mockUserRepo.On("FindByTelegramWithAuths", ctx, telegram).Return(user, nil)

	req := &dto.LoginRequest{Telegram: "@telegram", Password: "abracadabra"}
	code, _, _, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)

	mockUserRepo.AssertExpectations(t)
}
