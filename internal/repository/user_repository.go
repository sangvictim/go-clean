package repository

import (
	"go-clean/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	UserRepository *Repository[model.User]
	Log            *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) Search(db *gorm.DB, request *model.UserSearchRequest) ([]model.User, error) {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
