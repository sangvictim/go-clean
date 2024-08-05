package userUsecase

import (
	"context"
	userModel "go-clean/domain/user/model"
	userRepository "go-clean/domain/user/repository"
	apiResponse "go-clean/utils/api_response"
	"go-clean/utils/encrypt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *userRepository.UserRepository
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *userRepository.UserRepository) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUsecase) Search(ctx context.Context, request *userModel.UserSearchRequest) ([]userModel.UserResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validation user search")
		return nil, 0, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}
	users, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to search users")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit search user")
		return nil, 0, err
	}

	response := make([]userModel.UserResponse, len(users))
	for i, user := range users {
		response[i] = userModel.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			TimeStamp: user.TimeStamp,
		}
	}

	return response, total, nil
}

func (c *UserUsecase) FindById(ctx context.Context, request *userModel.Id) (*userModel.UserResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(userModel.User)

	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.WithError(err).Info("error getting user")
		return nil, echo.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit get user")
		return nil, echo.ErrInternalServerError
	}

	return &userModel.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		TimeStamp: user.TimeStamp,
	}, nil
}

func (c *UserUsecase) Create(ctx context.Context, request *userModel.UserCreate) (*userModel.UserResponse, error) {
	tx := c.DB.Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validation user create")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}

	checkUser := c.UserRepository.CheckByEmail(tx, request)
	if checkUser {
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "Email already exist",
		})
	}

	hasPassword, _ := encrypt.Brypt(request.Password)
	user := &userModel.User{
		UserEntity: userModel.UserEntity{
			Name:     request.Name,
			Email:    request.Email,
			Password: hasPassword,
		},
		TimeStamp: userModel.TimeStamp{
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("error creating user")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit create user")
		return nil, echo.ErrInternalServerError
	}

	return &userModel.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		TimeStamp: user.TimeStamp,
	}, nil
}

func (c *UserUsecase) Update(ctx context.Context, request *userModel.UserUpdate, id int) (*userModel.UserResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(userModel.User)
	c.UserRepository.FindById(tx, user, id)
	if user.ID == 0 {
		return nil, echo.NewHTTPError(404, apiResponse.Response{
			Message: "User not Found",
		})
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Info("error validation user update")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}
	hasPassword, _ := encrypt.Brypt(request.Password)
	user = &userModel.User{
		Id: user.Id,
		UserEntity: userModel.UserEntity{
			Name:     request.Name,
			Email:    request.Email,
			Password: hasPassword,
		},
		TimeStamp: userModel.TimeStamp{
			CreatedAt: user.CreatedAt,
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	if err := c.UserRepository.Update(tx, user, id); err != nil {
		c.Log.WithError(err).Error("error updating user")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit update user")
		return nil, echo.ErrInternalServerError
	}

	return &userModel.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		TimeStamp: user.TimeStamp,
	}, nil
}

func (c *UserUsecase) Delete(ctx context.Context, id int) error {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(userModel.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.WithError(err).Info("error getting user")
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
