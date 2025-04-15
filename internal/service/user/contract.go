package user

import "github.com/developeerz/restorio-auth/internal/repository/postgres/models"

type Repository interface {
	SetUserAuth(userAuth *models.UserAuth) error
	CreateUser(user *models.User) error
	FindByTelegram(telegram string) (*models.User, error)
	FindByTelegramWithAuths(telegram string) (*models.UserWithAuths, error)
}

type Cache interface {
	PutUser(telegram string, userJSON []byte) error
	PutVerificationCode(telegram string, code int) error
	GetUser(telegram string) ([]byte, error)
	GetVerificationCode(telegram string) (int, error)
}
