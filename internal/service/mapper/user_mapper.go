package mapper

import (
	"strings"
	"time"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUpToUser(signUp *dto.SignUpRequest) (*models.User, error) {
	signUp.Telegram, _ = strings.CutPrefix(signUp.Telegram, "@")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUp.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Firstname:        signUp.Firstname,
		Lastname:         signUp.Lastname,
		Telegram:         signUp.Telegram,
		Password:         string(hashedPassword),
		RegistrationDate: time.Now(),
	}, nil
}

func VerificationToUserCode(v *dto.VerificationRequest) (*models.UserCode, error) {
	v.Telegram, _ = strings.CutPrefix(v.Telegram, "@")

	return &models.UserCode{
		Telegram: v.Telegram,
		Code:     v.Code,
	}, nil
}

func UserAuthToIDAndAuth(userAuths []models.UserAuth) (int64, []string) {
	id := userAuths[0].UserID

	auths := make([]string, len(userAuths))
	for i, v := range userAuths {
		auths[i] = string(v.AuthID)
	}

	return id, auths
}

func JwtToAccess(jwts *jwt.Jwt) *dto.JwtAccess {
	return &dto.JwtAccess{
		Access: jwts.Access,
	}
}

func AuthsToStrings(auths []models.Authority) []string {
	strAuths := make([]string, len(auths))
	for i, v := range auths {
		strAuths[i] = string(v.ID)
	}

	return strAuths
}
