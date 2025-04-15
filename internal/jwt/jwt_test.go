package jwt_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/jwt"
	jwtpkg "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJwt(t *testing.T) {
	t.Parallel()

	token, err := jwt.NewJwt(1, []string{"USER", "ADMIN"})
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Access)
	assert.NotEmpty(t, token.Refresh)
}

func TestParseRefresh(t *testing.T) {
	t.Parallel()

	expectedtelegramID := int64(1)
	token, err := jwt.NewJwt(expectedtelegramID, []string{"USER", "ADMIN"})
	assert.NoError(t, err)

	userID, err := jwt.ParseRefresh(token.Refresh)
	assert.NoError(t, err)
	assert.Equal(t, expectedtelegramID, userID)
}

func TestGetAccess(t *testing.T) {
	t.Parallel()

	expectedtelegramID := int64(1)
	roles := []string{"USER", "ADMIN"}
	token, err := jwt.NewJwt(expectedtelegramID, roles)
	assert.NoError(t, err)

	strtelegramID, roleStrings, err := jwt.GetAccess(token.Access)
	assert.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(expectedtelegramID, 10), strtelegramID)
	assert.Equal(t, roles, roleStrings)
}

func TestGetValidToken(t *testing.T) {
	t.Parallel()

	claims := &jwtpkg.MapClaims{
		"sub":   1,
		"roles": []string{"USER"},
		"iat":   jwtpkg.NewNumericDate(time.Now()),
		"exp":   jwtpkg.NewNumericDate(time.Now().Add(-1 * time.Second * jwt.AccessMaxAge)),
	}

	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	access, err := token.SignedString([]byte(config.ConfigService.Access))
	assert.NoError(t, err)
	assert.NotEmpty(t, access)

	_, _, err = jwt.GetAccess(access)
	assert.Error(t, err)
}
