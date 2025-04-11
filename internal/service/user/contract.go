package user

import "github.com/developeerz/restorio-auth/internal/repository/models"

type Repository interface {
	CreateUserAuth(userAuth *models.UserAuth) error
	CreateVerificationCode(userCode *models.UserCode) error
	DeleteVerificationCode(userCode *models.UserCode) error
	CheckVerificationCode(userCode *models.UserCode) (int64, error)
	SetUserAuth(userAuth *models.UserAuth) error
	CreateUser(user *models.User) error
	SaveUser(user *models.User) error
	FindByTelegram(telegram string) (*models.User, error)
	FindByTelegramWithAuths(telegram string) (*models.UserWithAuths, error)
}
