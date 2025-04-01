package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/developeerz/restorio-auth/config"
	"github.com/golang-jwt/jwt/v5"
)

// seconds
const (
	AccessMaxAge  = 60 * 60
	RefreshMaxAge = 60 * 60 * 24 * 7
)

type Jwt struct {
	Access  string
	Refresh string
}

func genAccessToken(userId string, auths []string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub":   userId,
		"roles": auths,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Second * AccessMaxAge)),
	}
}

func genRefreshToken(userId string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": userId,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Second * RefreshMaxAge)),
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
	token, err := getValidToken(refreshToken, config.ConfigService.Refresh)

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

func GetAccess(accessToken string) (string, []string, error) {
	token, err := getValidToken(accessToken, config.ConfigService.Access)
	if err != nil {
		return "", nil, errors.New("invalid token")
	}

	strUserId, err := token.Claims.GetSubject()
	if err != nil {
		return "", nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.New("invalid claims")
	}

	roles, ok := claims["roles"].([]interface{})
	if !ok {
		return "", nil, errors.New("invalid roles")
	}

	var roleStrings []string
	for _, role := range roles {
		roleStrings = append(roleStrings, role.(string))
	}

	return strUserId, roleStrings, nil
}

func getValidToken(token string, key string) (*jwt.Token, error) {
	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !jwt.Valid {
		return nil, fmt.Errorf("Not valid token")
	}

	if expired(jwt) {
		return nil, fmt.Errorf("Token is expired")
	}

	return jwt, nil
}

func expired(token *jwt.Token) bool {
	date, err := token.Claims.GetExpirationTime()
	if err != nil {
		return true
	}

	return date.Time.Compare(time.Now()) == -1
}
