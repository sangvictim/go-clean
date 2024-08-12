package user

import (
	"go-clean/domain/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	repository.Repository[User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, entity *User, id any) error {
	return db.Where("email = ?", id).First(entity).Error
}

func (r *UserRepository) Search(db *gorm.DB, request *UserSearchRequest) ([]User, int64, error) {
	var users []User
	if err := db.Scopes(r.FilterUser(request)).Limit(request.Limit).Offset((request.Page - 1) * request.Limit).Order(
		clause.OrderByColumn{
			Column: clause.Column{Name: request.OrderBy},
			Desc:   request.OrderDirection == "desc",
		},
	).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FilterUser(request *UserSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if searchName := request.Search; searchName != "" {
			searchName = "%" + searchName + "%"
			tx = tx.Where("name LIKE ?", searchName).Or("email LIKE ?", searchName)
		}

		return tx
	}
}

func (r *UserRepository) IsEmail(db *gorm.DB, request *User) bool {
	var user User

	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return false
	}

	return true
}
