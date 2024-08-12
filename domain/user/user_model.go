package user

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement; column:id; " json:"id,omitempty" example:"1" `
	Name      string    `gorm:"column:name; type:varchar(255)" json:"name" validate:"required" example:"John Doe"`
	Email     string    `gorm:"column:email; type:varchar(255)" json:"email" validate:"required,email" example:"z6Ls1@example.com"`
	Password  string    `gorm:"column:password; type:varchar(255)" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" autoCreateTime:"true" example:"2024-08-12 02:45:26.704606+00"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" autoCreateTime:"true" example:"2024-08-12 02:45:26.704606+00"`
}

type UserCreate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type UserDetail struct {
	User
	Password string `json:"-"`
}

func UserToResponse(user *User) *User {
	return &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
