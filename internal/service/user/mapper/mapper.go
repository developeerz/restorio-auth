package mapper

import (
	"strings"
	"time"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	redis "github.com/developeerz/restorio-auth/pkg/repository/redis"
	"golang.org/x/crypto/bcrypt"
)

func SignUpToUser(signUp *dto.SignUpRequest) *redis.User {
	signUp.Telegram, _ = strings.CutPrefix(signUp.Telegram, "@")

	return &redis.User{
		Firstname: signUp.Firstname,
		Lastname:  signUp.Lastname,
		Telegram:  signUp.Telegram,
		Password:  signUp.Password,
	}
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

func UserToUser(user *redis.User) (*models.User, error) {
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
