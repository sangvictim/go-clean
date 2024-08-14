package user

import (
	"time"
)

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id,omitempty" example:"1" `
	Name      string    `gorm:"not null; column:name; type:varchar(100)" json:"name" validate:"required" example:"John Doe"`
	Email     string    `gorm:"uniqueIndex;not null;column:email; type:varchar(100)" json:"email" validate:"required,email" example:"z6Ls1@example.com"`
	Password  string    `gorm:"not null;column:password; type:varchar(100)" json:"password,omitempty"`
	Avatar    *string   `gorm:"column:avatar; type:varchar(255)" json:"avatar,omitempty" example:"https://example.com/public/avatar/avatar.png"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" autoCreateTime:"true" example:"2024-08-12 02:45:26.704606+00"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" autoCreateTime:"true" example:"2024-08-12 02:45:26.704606+00"`
}

type UserCreate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Avatar   string `json:"avatar"`
}

type UserUpdate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type UserDetail struct {
	User
	Password string `json:"-"`
}

type UserSearchRequest struct {
	Search         string `json:"search"`
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	OrderBy        string `json:"orderBy"`
	OrderDirection string `json:"orderDirection"`
}

func UserToResponse(user *User) *User {
	return &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
