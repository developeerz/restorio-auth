package user

import (
	"context"

	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
)

type Repository interface {
	CreateUserAuth(userAuth *models.UserAuth) error
	CreateUser(user *models.User) error
	FindByTelegram(telegram string) (*models.User, error)
	FindByTelegramWithAuths(telegram string) (*models.UserWithAuths, error)
}

type Cache interface {
	PutUser(ctx context.Context, telegram string, userJSON []byte) error
	PutVerificationCode(ctx context.Context, telegram string, code int) error
	GetUser(ctx context.Context, telegram string) ([]byte, error)
	GetVerificationCode(ctx context.Context, telegram string) (int, error)
}
