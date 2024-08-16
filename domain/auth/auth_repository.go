package auth

import (
	"go-clean/domain/user"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository struct {
	Log *logrus.Logger
}

func NewAuthRepository(log *logrus.Logger) *AuthRepository {
	return &AuthRepository{
		Log: log,
	}
}

func (r *AuthRepository) Register(db *gorm.DB, entity *Register) error {
	user := user.User{
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
	}
	return db.Create(&user).Error
}

func (r *AuthRepository) Login(db *gorm.DB, entity *user.User) (user.User, error) {
	var user user.User
	if err := db.Model(&user).Where("email = ?", entity.Email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthRepository) IsEmail(db *gorm.DB, request *Register) bool {
	var user user.User

	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return false
	}

	return true
}

func (r *AuthRepository) RefreshToken(db *gorm.DB, refreshToken *AccessToken) error {
	return db.Save(refreshToken).Error
}

func (r *AuthRepository) DeviceToken(db *gorm.DB, deviceToken *DeviceToken) error {
	return db.Save(deviceToken).Error
}
