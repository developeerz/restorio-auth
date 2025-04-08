package models

import "time"

type User struct {
	ID               int64  `gorm:"primaryKey"`
	Firstname        string `gorm:"size:63"`
	Lastname         string `gorm:"size:63"`
	Telegram         string `gorm:"unique,size:63"`
	TelegramID       int64  `gorm:"default:null"`
	Password         string `gorm:"size:255"`
	RegistrationDate time.Time
}

type UserCode struct {
	Telegram string `gorm:"primaryKey,unique,size:63"`
	Code     int
}

type UserWithAuths struct {
	ID       int64       `gorm:"primaryKey"`
	Telegram string      `gorm:"unique,size:63"`
	Password string      `gorm:"size:255"`
	Auths    []Authority `gorm:"many2many:user_auths;joinForeignKey:UserID;joinReferences:AuthID"`
}
