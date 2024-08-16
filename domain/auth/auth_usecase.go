package auth

import (
	"context"
	"go-clean/domain/user"
	"go-clean/utils/encrypt"
	"net/http"
	"sync"
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
}

func NewAuthUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, authRepository *AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		AuthRepository: authRepository,
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

	var accessToken, refreshToken string
	var timeExp, timeRefExp time.Time

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		accessToken, timeExp, _ = c.generateToken(res)
	}()

	go func() {
		defer wg.Done()
		refreshToken, timeRefExp, _ = c.generateRefreshToken(res)
	}()

	go func() {
		defer wg.Done()
		requestAccessToken := &AccessToken{
			UserId:       res.Id,
			RefreshToken: refreshToken,
			ExpiredAt:    timeRefExp,
		}
		c.createRefreshToken(tx, requestAccessToken)
	}()

	go func() {
		defer wg.Done()
		requestDevicetoken := &DeviceToken{
			UserId:      res.Id,
			DeviceID:    request.DeviceID,
			DeviceType:  request.DeviceType,
			UserAgent:   request.UserAgent,
			LastLoginAt: time.Now(),
		}
		c.createDeviceToken(tx, requestDevicetoken)
	}()

	wg.Wait()

	if err := tx.Commit().Error; err != nil {
		return nil, echo.ErrInternalServerError
	}

	return &LoginResponse{
		Id:        res.Id,
		Name:      res.Name,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		AccessToken: DetailToken{
			Token:       accessToken,
			TokenExpiry: timeExp,
		},
		RefreshToken: DetailToken{
			Token:       refreshToken,
			TokenExpiry: timeRefExp,
		},
	}, nil
}

func (c *AuthUsecase) Register(ctx context.Context, request *Register) error {
	tx := c.DB.Begin()

	checkUser := c.AuthRepository.IsEmail(tx, request)
	if checkUser {
		return echo.NewHTTPError(http.StatusConflict, "email already exist")
	}

	hasPassword, _ := encrypt.Brypt(request.Password)
	user := &Register{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasPassword,
	}

	if err := c.AuthRepository.Register(tx, user); err != nil {
		c.Log.WithError(err).Error("error register user")
		return echo.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit register user")
		return echo.ErrInternalServerError
	}

	return nil
}

func (c *AuthUsecase) generateToken(res user.User) (string, time.Time, error) {
	// masa aktif token 1 hari
	timeExp := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    res.Id,
		"email": res.Email,
		"exp":   timeExp,
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.key")))

	return tokenString, timeExp, err
}

func (c *AuthUsecase) generateRefreshToken(res user.User) (string, time.Time, error) {
	// masa aktif token 30 hari
	timeExp := time.Now().Add(time.Hour * 24 * 30)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    res.Id,
		"email": res.Email,
		"exp":   timeExp,
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.key")))

	return tokenString, timeExp, err
}

func (c *AuthUsecase) createRefreshToken(tx *gorm.DB, request *AccessToken) error {
	if err := c.AuthRepository.RefreshToken(tx, request); err != nil {
		return err
	}
	return nil
}

func (c *AuthUsecase) createDeviceToken(tx *gorm.DB, request *DeviceToken) error {
	if err := c.AuthRepository.DeviceToken(tx, request); err != nil {
		return err
	}
	return nil
}

func (c *AuthUsecase) Logout(ctx context.Context, token *string) error {
	tx := c.DB.Begin()

	// if err := c.AccessToken.Delete(tx, token); err != nil {
	// 	c.Log.WithError(err).Error("error delete token")
	// 	return echo.ErrInternalServerError
	// }

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit register user")
		return echo.ErrInternalServerError
	}

	return nil
}
