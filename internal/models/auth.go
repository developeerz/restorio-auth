package models

type Auth string

const (
	ADMIN Auth = "ADMIN"
	USER  Auth = "USER"
)

type Authority struct {
	ID          Auth   `gorm:"size:15"`
	Description string `gorm:"size:255"`
}

type UserAuth struct {
	UserId int64 `gorm:"primaryKey"`
	AuthId Auth  `gorm:"primaryKey"`
}
