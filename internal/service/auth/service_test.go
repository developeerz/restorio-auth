package auth_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/models"
	"github.com/developeerz/restorio-auth/internal/service/auth"
	"github.com/developeerz/restorio-auth/test/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRefreshSuccess(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userID := int64(1)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userID, auths)
	assert.NoError(t, err)

	userAuths := []models.UserAuth{{UserID: userID, AuthID: models.USER}}

	mockAuthRepo.On("GetUserAuths", userID).Return(userAuths, nil)

	_, _, err = service.Refresh(jwts.Refresh)
	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresParseError(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userID := int64(1)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userID, auths)
	assert.NoError(t, err)

	_, _, err = service.Refresh(fmt.Sprint(jwts.Refresh, "error"))

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestRefresGetUserAuthsError(t *testing.T) {
	t.Parallel()

	mockAuthRepo := mocks.NewRepository(t)
	service := auth.NewService(mockAuthRepo)

	userID := int64(1)
	auths := []string{string(models.USER)}
	jwts, err := jwt.NewJwt(userID, auths)
	assert.NoError(t, err)

	mockAuthRepo.On("GetUserAuths", userID).Return(nil, errors.New("Get user auths error"))

	_, _, err = service.Refresh(jwts.Refresh)

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}
