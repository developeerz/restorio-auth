package auth_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/models"
	"github.com/developeerz/restorio-auth/internal/service/auth"
	"github.com/developeerz/restorio-auth/internal/service/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRefreshSuccess(t *testing.T) {
	mockAuthRepo := mocks.NewAuthRepository(t)
	service := auth.NewAuthService(mockAuthRepo)

	userId := int64(1)
	auths := []string{string(models.USER)}
	jwts, _ := jwt.NewJwt(int64(userId), auths)
	userAuths := []models.UserAuth{{UserId: userId, AuthId: models.USER}}

	mockAuthRepo.On("GetUserAuths", userId).Return(userAuths, nil)

	_, _, err := service.Refresh(jwts.Refresh)

	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresParseError(t *testing.T) {
	mockAuthRepo := mocks.NewAuthRepository(t)
	service := auth.NewAuthService(mockAuthRepo)

	userId := int64(1)
	auths := []string{string(models.USER)}
	jwts, _ := jwt.NewJwt(int64(userId), auths)

	_, _, err := service.Refresh(fmt.Sprint(jwts.Refresh, "error"))

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresGetUserAuthsError(t *testing.T) {
	mockAuthRepo := mocks.NewAuthRepository(t)
	service := auth.NewAuthService(mockAuthRepo)

	userId := int64(1)
	auths := []string{string(models.USER)}
	jwts, _ := jwt.NewJwt(int64(userId), auths)

	mockAuthRepo.On("GetUserAuths", userId).Return(nil, errors.New("Get user auths error"))

	_, _, err := service.Refresh(jwts.Refresh)

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}
