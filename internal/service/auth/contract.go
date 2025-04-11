package auth

import "github.com/developeerz/restorio-auth/internal/repository/models"

type Repository interface {
	GetUserAuths(userID int64) ([]models.UserAuth, error)
}
