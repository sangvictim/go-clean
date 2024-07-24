package usecase

import (
	"context"
	"go-clean/internal/model"
	"go-clean/internal/repository"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUsecase) Search(ctx context.Context, request *model.UserSearchRequest) ([]model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	users, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get users")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, err
	}

	response := make([]model.UserResponse, len(users))
	for i, user := range users {
		response[i] = *model.UserToResponse(&user)
	}

	return response, nil

}
