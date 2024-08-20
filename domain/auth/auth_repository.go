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
	if err := db.Model(&user).Where("email = ?", entity.Email).Find(&user).Error; err != nil {
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

func (r *AuthRepository) FindByDeviceId(db *gorm.DB, deviceToken *DeviceToken) (DeviceToken, error) {
	var devices DeviceToken
	if err := db.Where("device_id = ?", deviceToken.DeviceID).Find(&devices).Error; err != nil {
		return devices, err
	}
	return devices, nil
}

func (r *AuthRepository) UpdateDevice(db *gorm.DB, deviceToken *DeviceToken) error {
	return db.Where("device_id = ?", deviceToken.DeviceID).Omit("created_at").Updates(deviceToken).Error
}

func (r *AuthRepository) DeleteToken(db *gorm.DB, refreshToken string) error {
	var accessToken AccessToken

	return db.Model(&accessToken).Where("refresh_token = ?", refreshToken).Delete(&accessToken).Error
}

func (r *AuthRepository) DeleteDevice(db *gorm.DB, device string) error {
	var devices DeviceToken
	return db.Model(&devices).Where("device_id = ?", device).Delete(&devices).Error
}
