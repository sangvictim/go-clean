package personalAccessToken

import (
	"go-clean/domain/user"
	"time"
)

type PersonalAccessToken struct {
	Id          uint `gorm:"primaryKey;autoIncrement;column:id"`
	UserId      uint `gorm:"column:user_id; type:bigint"`
	User        user.User
	AccessToken string    `gorm:"column:access_token; type:varchar(255)"`
	IP          string    `gorm:"column:ip; type:varchar(100)"`
	UserAgent   string    `gorm:"column:user_agent; type:varchar(100)"`
	ExpiredAt   time.Time `gorm:"column:expired_at"`
	CreatedAt   time.Time `gorm:"column:created_at" autoCreateTime:"true"`
}
