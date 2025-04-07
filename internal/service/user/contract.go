package user

import "github.com/developeerz/restorio-auth/internal/models"

type Repository interface {
	CreateUserAuth(userAuth *models.UserAuth) error
	GetUserAuths(userID int64) ([]models.UserAuth, error)
	CreateVerificationCode(userCode *models.UserCode) error
	DeleteVerificationCode(userCode *models.UserCode) error
	CheckVerificationCode(userCode *models.UserCode) (int64, error)
	SetUserAuth(userAuth *models.UserAuth) error
	CreateUser(user *models.User) error
	SaveUser(user *models.User) error
	VerifyUser(userID int64) error
	FindByTelegram(telegram string) (*models.User, error)
}
