package auth

import (
	"context"
	"go-clean/domain/user"
	"go-clean/pkg"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	AuthRepository *AuthRepository
}

func NewAuthService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, authRepository *AuthRepository) *AuthService {
	return &AuthService{
		DB:             db,
		Log:            log,
		Validate:       validate,
		AuthRepository: authRepository,
	}
}

func (c *AuthService) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
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

	encryptService := pkg.NewBcryptService()
	ok := encryptService.CompareHashBrypt(request.Password, res.Password)
	if !ok {
		return nil, echo.ErrUnauthorized
	}

	var accessToken, refreshToken string
	timeExp := time.Now().Add(time.Hour * 24)
	timeRefExp := time.Now().Add(time.Hour * 24 * 30)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	go func() {
		defer wg.Done()
		token, _ := c.generateJWTToken(res, timeExp)
		mu.Lock()
		accessToken = token
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		token, _ := c.generateJWTToken(res, timeRefExp)
		mu.Lock()
		refreshToken = token
		mu.Unlock()
	}()

	wg.Wait()

	requestAccessToken := &AccessToken{
		UserId:       res.Id,
		RefreshToken: refreshToken,
		ExpiredAt:    timeRefExp,
	}
	c.createRefreshToken(tx, requestAccessToken)

	requestDevicetoken := &DeviceToken{
		UserId:      res.Id,
		DeviceID:    request.DeviceID,
		DeviceType:  request.DeviceType,
		UserAgent:   request.UserAgent,
		LastLoginAt: time.Now(),
	}
	c.createDeviceToken(tx, requestDevicetoken)

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

func (c *AuthService) Register(ctx context.Context, request *Register) error {
	tx := c.DB.Begin()

	checkUser := c.AuthRepository.IsEmail(tx, request)
	if checkUser {
		return echo.NewHTTPError(http.StatusConflict, "email already exist")
	}

	encryptService := pkg.NewBcryptService()
	hasPassword, _ := encryptService.Bcrypt(request.Password)
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

func (c *AuthService) generateJWTToken(res user.User, time time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    res.Id,
		"email": res.Email,
		"exp":   time.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	return tokenString, err
}

func (c *AuthService) createRefreshToken(tx *gorm.DB, request *AccessToken) error {
	if err := c.AuthRepository.RefreshToken(tx, request); err != nil {
		return err
	}
	return nil
}

func (c *AuthService) createDeviceToken(tx *gorm.DB, request *DeviceToken) error {
	deviceToken, _ := c.AuthRepository.FindByDeviceId(tx, request)

	if deviceToken.Id != 0 {
		deviceToken.LastLoginAt = time.Now()
		if err := c.AuthRepository.UpdateDevice(tx, &deviceToken); err != nil {
			return err
		}
	} else {
		if err := c.AuthRepository.DeviceToken(tx, request); err != nil {
			return err
		}
	}
	return nil
}

func (c *AuthService) Logout(ctx context.Context, RefreshToken string, deviceToken string) error {
	tx := c.DB.Begin()
	defer tx.Rollback()

	if err := c.AuthRepository.DeleteToken(tx, RefreshToken); err != nil {
		return err
	}

	if err := c.AuthRepository.DeleteDevice(tx, deviceToken); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error commit register user")
		return echo.ErrInternalServerError
	}

	return nil
}
