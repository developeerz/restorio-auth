package jwt

import (
	"strconv"
	"time"

	"github.com/developeerz/restorio-auth/config"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func genAccessToken(userId int64, auths []string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": strconv.FormatInt(userId, 10),
		"aud": auths,
		"iss": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
}

func genRefreshToken(userId int64) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": strconv.FormatInt(userId, 10),
		"iss": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
}

func NewJwt(userId int64, auths []string) (*Jwt, error) {
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, genAccessToken(userId, auths))
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, genRefreshToken(userId))

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
