package postgres

import (
	"context"

	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByTelegram(ctx context.Context, telegram string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("telegram = ?", telegram).First(&user)

	return &user, result.Error
}

func (r *UserRepository) FindByTelegramWithAuths(ctx context.Context, telegram string) (*models.UserWithAuths, error) {
	var user models.UserWithAuths
	result := r.db.WithContext(ctx).
		Table("users u").
		Preload("Auths").
		Where("telegram = ?", telegram).
		First(&user)

	return &user, result.Error
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetUserAuths(ctx context.Context, telegramID int64) ([]models.UserAuth, error) {
	var userAuths []models.UserAuth
	result := r.db.WithContext(ctx).
		Table("user_auths").
		Where("user_telegram_id = ?", telegramID).
		Find(&userAuths)

	return userAuths, result.Error
}

func (r *UserRepository) CreateUserAuth(ctx context.Context, userAuth *models.UserAuth) error {
	return r.db.WithContext(ctx).Table("user_auths").Create(userAuth).Error
}
