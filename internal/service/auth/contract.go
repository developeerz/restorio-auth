package auth

import "github.com/developeerz/restorio-auth/internal/repository/postgres/models"

type Repository interface {
	GetUserAuths(telegramID int64) ([]models.UserAuth, error)
}
