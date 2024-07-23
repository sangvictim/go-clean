package models

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
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
