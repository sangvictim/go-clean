package model

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement; column:id; type:bigint" json:"id"`
	Name      string    `gorm:"column:name; type:varchar(255)" json:"name"`
	Email     string    `gorm:"column:email; type:varchar(255)" json:"email"`
	Password  string    `gorm:"column:password; type:varchar(255)" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"lQwLd@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type UserSearchRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1,max=100"`
}

func UserToResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
