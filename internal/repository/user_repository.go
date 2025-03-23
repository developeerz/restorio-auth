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
	FindVerificationCodeByUserID(userId int64, code int) error
	CheckVerificationCode(telegram string, code int) (int64, error)
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
	result := r.db.Save(user)
	return result.Error
}

func (r *userRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) VerifyUser(userId int64) error {
	result := r.db.Table("users").Where("id = ?", userId).Update("is_verified", true)
	return result.Error
}

func (r *userRepository) GetUserAuths(userId int64) ([]models.UserAuth, error) {
	var userAuths []models.UserAuth
	result := r.db.Table("user_auths").Where("user_id = ?", userId).Find(&userAuths)
	return userAuths, result.Error
}

func (r *userRepository) CreateUserAuth(userAuth *models.UserAuth) error {
	result := r.db.Create(userAuth)
	return result.Error
}

func (r *userRepository) CreateVerificationCode(userCode *models.UserCode) error {
	result := r.db.Create(userCode)
	return result.Error
}

func (r *userRepository) DeleteVerificationCode(userCode *models.UserCode) error {
	result := r.db.Delete(userCode)
	return result.Error
}

func (r *userRepository) FindVerificationCodeByUserID(userId int64, code int) error {
	result := r.db.Table("user_codes").Where("user_id = ? and code = ?", userId, code)
	return result.Error
}

func (r *userRepository) CheckVerificationCode(telegram string, code int) (int64, error) {
	var userCode *models.UserCode
	result := r.db.Table("user_codes uc").
		Joins("JOIN users us ON us.id = uc.user_id").
		Where("us.telegram = ? AND uc.code = ?", telegram, code).
		First(&userCode)
	return userCode.UserId, result.Error
}

func (r *userRepository) SetUserAuth(userAuth *models.UserAuth) error {
	return r.db.Table("user_auths").Create(userAuth).Error
}
