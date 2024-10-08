package user

import (
	"context"
	"go-clean/pkg"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *UserRepository
}

func NewUserService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *UserRepository) *UserService {
	return &UserService{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserService) Search(ctx context.Context, request *UserSearchRequest) ([]User, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	users, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("failed to search users")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit search user")
		return nil, 0, err
	}

	response := make([]User, len(users))
	for i, user := range users {
		response[i] = *UserToResponse(&user)
	}

	return response, total, nil
}

func (c *UserService) FindById(ctx context.Context, id int) (*User, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)

	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		return nil, echo.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
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

func (c *UserService) Create(ctx context.Context, request *User) (*User, error) {
	tx := c.DB.Begin()

	checkUser := c.UserRepository.IsEmail(tx, request)
	if checkUser {
		return nil, echo.NewHTTPError(http.StatusConflict, "email already exist")
	}

	hasPassword, _ := pkg.NewBcryptService().Bcrypt(request.Password)
	user := &User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasPassword,
		Avatar:   request.Avatar,
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
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *UserService) Update(ctx context.Context, request *User, id int) (*UserDetail, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)
	c.UserRepository.FindById(tx, user, id)
	if user.Id == 0 {
		return nil, echo.ErrNotFound
	}

	hasPassword, _ := pkg.NewBcryptService().Bcrypt(request.Password)
	requestForm := &User{
		Id:        user.Id,
		Name:      request.Name,
		Email:     request.Email,
		Password:  hasPassword,
		Avatar:    user.Avatar,
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

func (c *UserService) Delete(ctx context.Context, id int) error {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := new(User)
	c.UserRepository.FindById(tx, user, id)
	if user.Id == 0 {
		return echo.ErrNotFound
	}

	if err := c.UserRepository.Delete(tx, user); err != nil {
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}
