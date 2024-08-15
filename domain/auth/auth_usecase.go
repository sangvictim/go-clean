package auth

import (
	"context"
	personalAccessToken "go-clean/domain/personalAccesstoken"
	"go-clean/domain/user"
	"go-clean/utils/encrypt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	AuthRepository *AuthRepository
	AccessToken    *personalAccessToken.PersonalAccessTokenRepository
}

func NewAuthUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, authRepository *AuthRepository, accessToken *personalAccessToken.PersonalAccessTokenRepository) *AuthUsecase {
	return &AuthUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		AuthRepository: authRepository,
		AccessToken:    accessToken,
	}
}

func (c *AuthUsecase) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	tx := c.DB.Begin()
	defer tx.Rollback()

	user := &user.User{
		Email:    request.Email,
		Password: request.Password,
	}

	res, err := c.AuthRepository.Login(tx, user)
	if err != nil {
		return nil, echo.ErrUnauthorized
	}

	if !encrypt.CompareHashBrypt(request.Password, res.Password) {
		return nil, echo.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    res.Id,
		"email": res.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return nil, err
	}

	c.createAccessToken(tx, res.Id, tokenString, request)

	if err := tx.Commit().Error; err != nil {
		return nil, echo.ErrInternalServerError
	}

	return &LoginResponse{
		Id:        res.Id,
		Name:      res.Name,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Token:     tokenString,
	}, nil
}

func (c *AuthUsecase) Register(ctx context.Context, request *Register) (*Register, error) {
	tx := c.DB.Begin()

	checkUser := c.AuthRepository.IsEmail(tx, request)
	if checkUser {
		return nil, echo.NewHTTPError(http.StatusConflict, "email already exist")
	}

	hasPassword, _ := encrypt.Brypt(request.Password)
	user := &Register{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasPassword,
	}

	if err := c.AuthRepository.Register(tx, user); err != nil {
		c.Log.WithError(err).Error("error register user")
		return nil, echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit register user")
		return nil, echo.ErrInternalServerError
	}

	return user, nil
}

func (c *AuthUsecase) createAccessToken(tx *gorm.DB, userId uint, tokenString string, request *LoginRequest) {
	requestToken := &personalAccessToken.PersonalAccessToken{
		UserId:      userId,
		AccessToken: tokenString,
		IP:          request.Ip,
		UserAgent:   request.UserAgent,
		ExpiredAt:   time.Now().Add(time.Hour * 24),
	}

	c.AccessToken.Create(tx, requestToken)
}

func (c *AuthUsecase) Logout(ctx context.Context, token *string) error {
	tx := c.DB.Begin()

	if err := c.AccessToken.Delete(tx, token); err != nil {
		c.Log.WithError(err).Error("error delete token")
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit register user")
		return echo.ErrInternalServerError
	}

	return nil
}
