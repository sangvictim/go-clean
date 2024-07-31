package usecase

import (
	"context"
	"go-clean/internal/model"
	"go-clean/internal/repository"
	apiResponse "go-clean/utils/api_response"
	"go-clean/utils/encrypt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewAuthUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *AuthUsecase) Register(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("validation request")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}

	hashPassword, _ := encrypt.Brypt(request.Password)

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("error registering user")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error registering user")
		return nil, echo.ErrInternalServerError
	}

	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
