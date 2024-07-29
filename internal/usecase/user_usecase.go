package usecase

import (
	"context"
	"go-clean/internal/model"
	"go-clean/internal/repository"
	apiResponse "go-clean/utils/api_response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func (c *UserUsecase) Search(ctx context.Context, request *model.UserSearchRequest) ([]model.UserResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("validation request")
		return nil, 0, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}
	users, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get users")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, err
	}

	response := make([]model.UserResponse, len(users))
	for i, user := range users {
		response[i] = *model.UserToResponse(&user)
	}

	return response, total, nil

}
