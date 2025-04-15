package models

import "time"

type Auth string

const (
	ADMIN Auth = "ADMIN"
	USER  Auth = "USER"
)

type User struct {
	TelegramID       int64  `gorm:"primaryKey"`
	Firstname        string `gorm:"size:63"`
	Lastname         string `gorm:"size:63"`
	Telegram         string `gorm:"unique,size:63"`
	Password         string `gorm:"size:255"`
	RegistrationDate time.Time
}

type UserWithAuths struct {
	TelegramID int64       `gorm:"primaryKey"`
	Telegram   string      `gorm:"unique,size:63"`
	Password   string      `gorm:"size:255"`
	Auths      []Authority `gorm:"many2many:user_auths;joinForeignKey:TelegramID;joinReferences:AuthID"`
}

type Authority struct {
	ID          Auth   `gorm:"size:15"`
	Description string `gorm:"size:255"`
}

type UserAuth struct {
	TelegramID int64 `gorm:"primaryKey"`
	AuthID     Auth  `gorm:"primaryKey"`
}

func (Authority) TableName() string {
	return "auths"
}
