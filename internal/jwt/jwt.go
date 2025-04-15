package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/developeerz/restorio-auth/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessMaxAge  = 60 * 60
	RefreshMaxAge = 60 * 60 * 24 * 7
)

type Jwt struct {
	Access  string
	Refresh string
}

func genAccessToken(telegramID string, auths []string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub":   telegramID,
		"roles": auths,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Second * AccessMaxAge)),
	}
}

func genRefreshToken(telegramID string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": telegramID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Second * RefreshMaxAge)),
	}
}

func NewJwt(telegramID int64, auths []string) (*Jwt, error) {
	strtelegramID := strconv.FormatInt(telegramID, 10)
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, genAccessToken(strtelegramID, auths))
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, genRefreshToken(strtelegramID))

	a, err := access.SignedString([]byte(config.ConfigService.Access))
	if err != nil {
		return nil, err
	}

	r, err := refresh.SignedString([]byte(config.ConfigService.Refresh))
	if err != nil {
		return nil, err
	}

	return &Jwt{
		Access:  a,
		Refresh: r,
	}, nil
}

func ParseRefresh(refreshToken string) (int64, error) {
	token, err := getValidToken(refreshToken, config.ConfigService.Refresh)
	if err != nil {
		return 0, err
	}

	strtelegramID, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	telegramID, err := strconv.ParseInt(strtelegramID, 10, 64)
	if err != nil {
		return 0, err
	}

	return telegramID, nil
}

func GetAccess(accessToken string) (string, []string, error) {
	token, err := getValidToken(accessToken, config.ConfigService.Access)
	if err != nil {
		return "", nil, fmt.Errorf("invalid token")
	}

	strtelegramID, err := token.Claims.GetSubject()
	if err != nil {
		return "", nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, fmt.Errorf("token.Claims: does not provide <jwt.MapClaims> value")
	}

	roles, ok := claims["roles"].([]interface{})
	if !ok {
		return "", nil, fmt.Errorf("claims[\"roles\"]: does not provide <[]interface{}> value")
	}

	roleStrings := make([]string, len(roles))

	for i, role := range roles {
		strRole, ok := role.(string)

		if !ok {
			return "", nil, fmt.Errorf("role: does not provide <string> value")
		}
		roleStrings[i] = strRole
	}

	return strtelegramID, roleStrings, nil
}

func getValidToken(token string, key string) (*jwt.Token, error) {
	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !jwt.Valid {
		return nil, fmt.Errorf("not valid token")
	}

	if expired(jwt) {
		return nil, fmt.Errorf("token is expired")
	}

	return jwt, nil
}

func expired(token *jwt.Token) bool {
	date, err := token.Claims.GetExpirationTime()
	if err != nil {
		return true
	}

	return date.Compare(time.Now()) == -1
}
