package repository

import (
	"go-clean/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	UserRepository *Repository[models.User]
	DB             *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
