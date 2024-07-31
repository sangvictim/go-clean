package usecase

import (
	"context"
	"go-clean/internal/model"
	"go-clean/internal/repository"
	apiResponse "go-clean/utils/api_response"
	"go-clean/utils/encrypt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

func (c *AuthUsecase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	tx := c.DB.Begin()
	var (
		key   string
		t     *jwt.Token
		token string
	)
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("validation request")
		return nil, echo.NewHTTPError(400, apiResponse.Response{
			Message: "validation request",
			Errors:  err.Error(),
		})
	}

	user, err := c.UserRepository.FindByEmail(tx, request)
	if err != nil {
		return nil, echo.NewHTTPError(404, apiResponse.Response{
			Message: "User not Found",
			Errors:  err.Error(),
		})
	}

	// Generate JWT token
	key = "sangvictim"
	t = jwt.New(jwt.SigningMethodHS256)
	claims := t.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token, err = t.SignedString([]byte(key))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, apiResponse.Response{
			Message: "Failed to generate token",
			Errors:  err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error login")
		return nil, echo.ErrInternalServerError
	}

	return &model.LoginResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
