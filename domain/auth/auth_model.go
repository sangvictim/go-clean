package auth

import (
	"go-clean/domain/user"
	"time"
)

type Register struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email" example:"super@mail.com"`
	Password  string `json:"password" validate:"required" example:"123"`
	Ip        string `json:"-"`
	UserAgent string `json:"-"`
}

type LoginResponse struct {
	Id           uint        `json:"id"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	AccessToken  DetailToken `json:"access_token"`
	RefreshToken DetailToken `json:"refresh_token"`
}

type DetailToken struct {
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"token_expiry"`
}

type AccessToken struct {
	Id           uint `gorm:"primaryKey;autoIncrement;column:id"`
	UserId       uint `gorm:"column:user_id; type:bigint"`
	User         user.User
	RefreshToken string    `gorm:"column:refresh_token; type:varchar(255)"`
	ExpiredAt    time.Time `gorm:"column:expired_at"`
	CreatedAt    time.Time `gorm:"column:created_at" autoCreateTime:"true"`
}

type DeviceToken struct {
	Id          uint `gorm:"primaryKey;autoIncrement;column:id"`
	UserId      uint `gorm:"column:user_id; type:bigint"`
	User        user.User
	DeviceID    string    `gorm:"column:device_id; type:varchar(255)"`
	DeviceName  string    `gorm:"column:device_name; type:varchar(100)"`
	FcmToken    string    `gorm:"column:fcm_token; type:varchar(255)"`
	IP          string    `gorm:"column:ip; type:varchar(100)"`
	UserAgent   string    `gorm:"column:user_agent; type:varchar(100)"`
	LastLoginAt time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time `gorm:"column:created_at" autoCreateTime:"true"`
}
