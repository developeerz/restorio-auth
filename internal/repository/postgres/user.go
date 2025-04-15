package postgres

import (
	"github.com/developeerz/restorio-auth/internal/repository/postgres/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByTelegram(telegram string) (*models.User, error) {
	var user models.User
	result := r.db.Where("telegram = ?", telegram).First(&user)

	return &user, result.Error
}

func (r *UserRepository) FindByTelegramWithAuths(telegram string) (*models.UserWithAuths, error) {
	var user models.UserWithAuths
	result := r.db.
		Table("users u").
		Preload("Auths").
		Where("telegram = ?", telegram).
		First(&user)

	return &user, result.Error
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserAuths(telegramID int64) ([]models.UserAuth, error) {
	var userAuths []models.UserAuth
	result := r.db.
		Table("user_auths").
		Where("user_telegram_id = ?", telegramID).
		Find(&userAuths)

	return userAuths, result.Error
}

func (r *UserRepository) SetUserAuth(userAuth *models.UserAuth) error {
	return r.db.Table("user_auths").Create(userAuth).Error
}
