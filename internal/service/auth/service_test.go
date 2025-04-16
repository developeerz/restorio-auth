package auth_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"github.com/developeerz/restorio-auth/internal/service/auth"
	mocks "github.com/developeerz/restorio-auth/test/mocks/auth"
	"github.com/stretchr/testify/assert"
)

func TestRefreshSuccess(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userTelegramID := int64(1111111)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userTelegramID, auths)
	assert.NoError(t, err)

	userAuths := []models.UserAuth{{UserTelegramID: userTelegramID, AuthID: models.USER}}

	mockAuthRepo.On("GetUserAuths", userTelegramID).Return(userAuths, nil)

	_, _, err = service.Refresh(jwts.Refresh)
	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresParseError(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userTelegramID := int64(111111)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userTelegramID, auths)
	assert.NoError(t, err)

	_, _, err = service.Refresh(fmt.Sprint(jwts.Refresh, "error"))

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresGetUserAuthsError(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userTelegramID := int64(111111)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userTelegramID, auths)
	assert.NoError(t, err)

	mockAuthRepo.On("GetUserAuths", userTelegramID).Return(nil, errors.New("Get user auths error"))

	_, _, err = service.Refresh(jwts.Refresh)

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}
