package repository

import (
	"github.com/developeerz/restorio-auth/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUserAuth(userAuth *models.UserAuth) error
	GetUserAuths(userId int64) ([]models.UserAuth, error)
	CreateVerificationCode(userCode *models.UserCode) error
	DeleteVerificationCode(userCode *models.UserCode) error
	CheckVerificationCode(userCode *models.UserCode) (int64, error)
	SetUserAuth(userAuth *models.UserAuth) error
	CreateUser(user *models.User) error
	SaveUser(user *models.User) error
	VerifyUser(userId int64) error
	FindByTelegram(telegram string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRipository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByTelegram(telegram string) (*models.User, error) {
	var user models.User
	result := r.db.Where("telegram = ?", telegram).First(&user)
	return &user, result.Error
}

func (r *userRepository) SaveUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) VerifyUser(userId int64) error {
	return r.db.Table("users").Where("id = ?", userId).Update("verified", true).Error
}

func (r *userRepository) GetUserAuths(userId int64) ([]models.UserAuth, error) {
	var userAuths []models.UserAuth
	result := r.db.Table("user_auths").Where("user_id = ?", userId).Find(&userAuths)
	return userAuths, result.Error
}

func (r *userRepository) CreateUserAuth(userAuth *models.UserAuth) error {
	return r.db.Create(userAuth).Error
}

func (r *userRepository) CreateVerificationCode(userCode *models.UserCode) error {
	return r.db.Create(userCode).Error
}

func (r *userRepository) DeleteVerificationCode(userCode *models.UserCode) error {
	return r.db.Table("user_codes").Where("telegram = ?", userCode.Telegram).Delete(userCode).Error
}

func (r *userRepository) CheckVerificationCode(userCode *models.UserCode) (int64, error) {
	var user models.User
	result := r.db.Select("us.id").Table("users us").
		Joins("JOIN user_codes uc ON uc.telegram = us.telegram").
		Where("uc.telegram = ? AND uc.code = ?", userCode.Telegram, userCode.Code).
		First(&user)
	return user.ID, result.Error
}

func (r *userRepository) SetUserAuth(userAuth *models.UserAuth) error {
	return r.db.Table("user_auths").Create(userAuth).Error
}
