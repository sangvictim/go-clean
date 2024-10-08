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
	Email      string `json:"email" validate:"required,email" example:"super@mail.com"`
	Password   string `json:"password" validate:"required" example:"123"`
	UserAgent  string `json:"-"`
	DeviceID   string `json:"-"`
	DeviceType string `json:"-"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"123"`
}

type LoginResponse struct {
	Id           uint        `json:"id"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	AccessToken  DetailToken `json:"access_token"`
	RefreshToken DetailToken `json:"refresh_token"`
}

type DetailToken struct {
	Token       string `json:"token"`
	TokenExpiry string `json:"token_expiry"`
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
	DeviceType  string    `gorm:"column:device_type; type:varchar(100)"`
	FcmToken    *string   `gorm:"column:fcm_token; type:varchar(255)"`
	UserAgent   string    `gorm:"column:user_agent; type:varchar(100)"`
	LastLoginAt time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time `gorm:"column:created_at" autoCreateTime:"true"`
}
