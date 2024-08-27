package log

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DBHook struct {
	DB *gorm.DB
}

func (hook *DBHook) Fire(entry *logrus.Entry) error {
	createdAt := entry.Time
	level := entry.Level.String()
	message := entry.Message

	err := hook.DB.Create(&Log{
		LogEntity: LogEntity{Level: level, Message: message},
		TimeStamp: TimeStamp{CreatedAt: createdAt, UpdatedAt: createdAt},
	}).Error
	return err
}

func (hook *DBHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
