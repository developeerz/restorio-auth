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
	Verified         bool `gorm:"default:false"`
}

type UserCode struct {
	Telegram string `gorm:"primaryKey,unique,size:63"`
	Code     int
}
