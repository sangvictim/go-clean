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

func (c *UserUsecase) FindById(ctx context.Context, request *model.UserId) (*model.UserResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(model.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		return nil, echo.NewHTTPError(404, apiResponse.Response{
			Message: "User not Found",
			Errors:  err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting contact")
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

func (c *UserUsecase) Create(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error) {
	tx := c.DB.Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("validation request")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, echo.ErrInternalServerError
	}

	// return converter.UserToResponse(user), nil
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *UserUsecase) Update(ctx context.Context, request *model.UserRequest, id int) (*model.UserResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(model.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		return nil, echo.NewHTTPError(404, apiResponse.Response{
			Message: "User not Found",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("validation request")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}

	user.ID = id
	user.Name = request.Name
	user.Email = request.Email
	user.Password = request.Password

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating contact")
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

func (c *UserUsecase) Delete(ctx context.Context, id int) error {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(model.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		return echo.NewHTTPError(404, apiResponse.Response{
			Message: "User not Found",
			Errors:  err.Error(),
		})
	}

	if err := c.UserRepository.Delete(tx, user); err != nil {
		c.Log.WithError(err).Error("error deleting user")
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error deleting user")
		return echo.ErrInternalServerError
	}

	return nil
}
