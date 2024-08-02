package userRepository

import (
	userModel "go-clean/domain/user/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[userModel.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) Search(db *gorm.DB, request *userModel.UserSearchRequest) ([]userModel.User, int64, error) {
	var users []userModel.User
	if err := db.Scopes(r.FilterUser(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&userModel.User{}).Scopes(r.FilterUser(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FilterUser(request *userModel.UserSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if id := request.Id; id != 0 {
			tx = tx.Where("id = ?", id)
		}

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if request.Email != "" {
			tx = tx.Where("email = ?", request.Email)
		}
		return tx
	}
}

// TODO: make it pointer request is generic, i need to use pointer in deference to pointer
func (r *UserRepository) FindByEmail(db *gorm.DB, request *userModel.UserEntity) (userModel.User, error) {
	var user userModel.User
	var total int64 = 0

	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return user, err
	}

	if err := db.Model(&userModel.User{}).Where("email = ?", request.Email).Count(&total).Error; err != nil {
		return user, err
	}
	return user, nil
}
