package auth

import "github.com/developeerz/restorio-auth/internal/models"

type AuthRepository interface {
	GetUserAuths(userId int64) ([]models.UserAuth, error)
}
