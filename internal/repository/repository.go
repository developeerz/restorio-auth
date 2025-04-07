package repository

import (
	"github.com/developeerz/restorio-auth/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByTelegram(telegram string) (*models.User, error) {
	var user models.User
	result := r.db.Where("telegram = ?", telegram).First(&user)

	return &user, result.Error
}

func (r *Repository) SaveUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) VerifyUser(userID int64) error {
	return r.db.Table("users").Where("id = ?", userID).Update("verified", true).Error
}

func (r *Repository) GetUserAuths(userID int64) ([]models.UserAuth, error) {
	var userAuths []models.UserAuth
	result := r.db.Table("user_auths").Where("user_id = ?", userID).Find(&userAuths)

	return userAuths, result.Error
}

func (r *Repository) CreateUserAuth(userAuth *models.UserAuth) error {
	return r.db.Create(userAuth).Error
}

func (r *Repository) CreateVerificationCode(userCode *models.UserCode) error {
	return r.db.Create(userCode).Error
}

func (r *Repository) DeleteVerificationCode(userCode *models.UserCode) error {
	return r.db.Table("user_codes").Where("telegram = ?", userCode.Telegram).Delete(userCode).Error
}

func (r *Repository) CheckVerificationCode(userCode *models.UserCode) (int64, error) {
	var user models.User
	result := r.db.Select("us.id").Table("users us").
		Joins("JOIN user_codes uc ON uc.telegram = us.telegram").
		Where("uc.telegram = ? AND uc.code = ?", userCode.Telegram, userCode.Code).
		First(&user)

	return user.ID, result.Error
}

func (r *Repository) SetUserAuth(userAuth *models.UserAuth) error {
	return r.db.Table("user_auths").Create(userAuth).Error
}
