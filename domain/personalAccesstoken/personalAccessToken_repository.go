package personalAccessToken

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonalAccessTokenRepository struct {
	Log *logrus.Logger
}

func NewPersonalAccessTokenRepository(log *logrus.Logger) *PersonalAccessTokenRepository {
	return &PersonalAccessTokenRepository{
		Log: log,
	}
}

func (r *PersonalAccessTokenRepository) Create(db *gorm.DB, request *PersonalAccessToken) error {
	return db.Create(request).Error
}

func (c *PersonalAccessTokenRepository) Delete(db *gorm.DB, token *string) error {
	return db.Where("access_token = ?", *token).Delete(PersonalAccessToken{}).Error
}
