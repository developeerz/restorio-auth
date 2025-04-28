package auth

import (
	"context"

	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
)

type Repository interface {
	GetUserAuths(ctx context.Context, telegramID int64) ([]models.UserAuth, error)
}
