package jwt

import (
	"errors"
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

func genAccessToken(userID string, auths []string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub":   userID,
		"roles": auths,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Second * AccessMaxAge)),
	}
}

func genRefreshToken(userID string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": userID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Second * RefreshMaxAge)),
	}
}

func NewJwt(userID int64, auths []string) (*Jwt, error) {
	strUserID := strconv.FormatInt(userID, 10)
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, genAccessToken(strUserID, auths))
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, genRefreshToken(strUserID))

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

	strUserID, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func GetAccess(accessToken string) (string, []string, error) {
	token, err := getValidToken(accessToken, config.ConfigService.Access)
	if err != nil {
		return "", nil, errors.New("invalid token")
	}

	strUserID, err := token.Claims.GetSubject()
	if err != nil {
		return "", nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.New("token.Claims: does not provide <jwt.MapClaims> value")
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

	return strUserID, roleStrings, nil
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
