package repository

import (
	"go-clean/internal/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DBHook struct {
	DB *gorm.DB
}

func (hook *DBHook) Fire(entry *logrus.Entry) error {
	createdAt := entry.Time.Format(time.RFC3339)
	level := entry.Level.String()
	message := entry.Message

	// query := `INSERT INTO logs (level, message, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	// err := hook.DB.Exec(query, level, message, createdAt, createdAt).Error
	err := hook.DB.Create(&model.Log{Level: level, Message: message, CreatedAt: createdAt, UpdatedAt: createdAt}).Error
	return err
}

func (hook *DBHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
