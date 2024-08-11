package user

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement; column:id; " json:"id,omitempty"`
	Name      string    `gorm:"column:name; type:varchar(255)" json:"name" validate:"required"`
	Email     string    `gorm:"column:email; type:varchar(255)" json:"email" validate:"required,email"`
	Password  string    `gorm:"column:password; type:varchar(255)" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" autoCreateTime:"true" `
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" autoCreateTime:"true" `
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
