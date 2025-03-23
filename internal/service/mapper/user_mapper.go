package mapper

import (
	"time"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUpToUser(signUp *dto.SignUpRequest) (*models.User, error) {
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
		IsVerified:       false,
	}, nil
}

func UserAuthToIdAndAuth(userAuths []models.UserAuth) (int64, []string) {
	id := userAuths[0].UserId

	var auths []string
	for _, v := range userAuths {
		auths = append(auths, string(v.AuthId))
	}

	return id, auths
}
