package models

import "time"

type User struct {
	ID               int64  `gorm:"primaryKey"`
	Firstname        string `gorm:"size:63"`
	Lastname         string `gorm:"size:63"`
	Telegram         string `gorm:"unique,size:255"`
	Password         string `gorm:"size:255"`
	RegistrationDate time.Time
	IsVerified       bool `gorm:"default:false"`
}

type UserCode struct {
	UserId int64 `gorm:"primaryKey"`
	Code   int
}
