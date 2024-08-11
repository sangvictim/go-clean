package user

import (
	"context"
	"go-clean/utils/encrypt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *UserRepository
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *UserRepository) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

// func (c *UserUsecase) Search(ctx context.Context, request *UserSearchRequest) ([]UserResponse, int64, error) {
// 	tx := c.DB.WithContext(ctx).Begin()
// 	defer tx.Rollback()

// 	if err := c.Validate.Struct(request); err != nil {
// 		c.Log.WithError(err).Error("error validation user search")
// 		return nil, 0, err
// 	}
// 	users, total, err := c.UserRepository.Search(tx, request)
// 	if err != nil {
// 		c.Log.WithError(err).Error("failed to search users")
// 		return nil, 0, err
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		c.Log.WithError(err).Error("error commit search user")
// 		return nil, 0, err
// 	}

// 	response := make([]UserResponse, len(users))
// 	for i, user := range users {
// 		response[i] = UserResponse{
// 			Name:  user.Name,
// 			Email: user.Email,
// 		}
// 	}

// 	return response, total, nil
// }

func (c *UserUsecase) FindById(ctx context.Context, id int) (*User, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)

	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.WithError(err).Info("error getting user")
		return nil, echo.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit get user")
		return nil, echo.ErrInternalServerError
	}

	return &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *UserUsecase) Create(ctx context.Context, request *User) (*User, error) {
	tx := c.DB.Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validation user create")
		return nil, err
	}

	checkUser := c.UserRepository.IsEmail(tx, request)
	if checkUser {
		return nil, echo.NewHTTPError(http.StatusConflict, "email already exist")
	}

	hasPassword, _ := encrypt.Brypt(request.Password)
	user := &User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasPassword,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("error creating user")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit create user")
		return nil, echo.ErrInternalServerError
	}

	return &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *UserUsecase) Update(ctx context.Context, request *User, id int) (*UserDetail, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)
	c.UserRepository.FindById(tx, user, id)
	if user.Id == 0 {
		return nil, echo.ErrNotFound
	}

	if err := c.Validate.Struct(request); err != nil {
		validationError := err.(validator.ValidationErrors)
		for _, validation := range validationError {
			return nil, echo.NewHTTPError(http.StatusBadRequest, validation)
		}
	}
	hasPassword, _ := encrypt.Brypt(request.Password)
	requestForm := &User{
		Id:        user.Id,
		Name:      request.Name,
		Email:     request.Email,
		Password:  hasPassword,
		CreatedAt: user.CreatedAt,
	}

	if err := c.UserRepository.Update(tx, requestForm, id); err != nil {
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error(err.Error())
		return nil, echo.ErrInternalServerError
	}

	return &UserDetail{
		User: *UserToResponse(requestForm),
	}, nil
}

func (c *UserUsecase) Delete(ctx context.Context, id int) error {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.WithError(err).Info("error getting user")
		return echo.ErrNotFound
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
