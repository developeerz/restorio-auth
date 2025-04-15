package mapper

import (
	"strings"
	"time"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	redis_dto "github.com/developeerz/restorio-auth/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

func SignUpToUser(signUp *dto.SignUpRequest) *redis_dto.User {
	signUp.Telegram, _ = strings.CutPrefix(signUp.Telegram, "@")

	return &redis_dto.User{
		Firstname: signUp.Firstname,
		Lastname:  signUp.Lastname,
		Telegram:  signUp.Telegram,
		Password:  signUp.Password,
	}
}

func UserAuthToIDAndAuth(userAuths []models.UserAuth) (int64, []string) {
	id := userAuths[0].TelegramID

	auths := make([]string, len(userAuths))
	for i, v := range userAuths {
		auths[i] = string(v.AuthID)
	}

	return id, auths
}

func JwtToAccess(jwts *jwt.Jwt) *dto.JwtAccessResponse {
	return &dto.JwtAccessResponse{
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

func UserToUser(user *redis_dto.User) (*models.User, error) {
	user.Telegram, _ = strings.CutPrefix(user.Telegram, "@")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &models.User{
		TelegramID:       user.TelegramID,
		Firstname:        user.Firstname,
		Lastname:         user.Lastname,
		Telegram:         user.Telegram,
		Password:         string(hashedPassword),
		RegistrationDate: time.Now(),
	}, nil
}
