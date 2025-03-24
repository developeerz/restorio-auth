package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/developeerz/restorio-auth/config"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func genAccessToken(userId string, auths []string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": userId,
		"aud": auths,
		"iss": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
}

func genRefreshToken(userId string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": userId,
		"iss": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
}

func NewJwt(userId int64, auths []string) (*Jwt, error) {
	strUserId := strconv.FormatInt(userId, 10)
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, genAccessToken(strUserId, auths))
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, genRefreshToken(strUserId))

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
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ConfigService.Refresh), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("Not valid token")
	}

	if isExpired(token) {
		return 0, fmt.Errorf("Token is expired")
	}

	strUserId, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(strUserId, 10, 64)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func isExpired(token *jwt.Token) bool {
	date, err := token.Claims.GetExpirationTime()
	if err != nil {
		return true
	}
	return date.Time.Compare(time.Now()) == -1
}
