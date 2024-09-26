package config

import (
	"fmt"
	"go-clean/domain/auth"
	"go-clean/domain/log"
	"go-clean/domain/user"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(viper *viper.Viper, Log *logrus.Logger) *gorm.DB {
	DB_HOST := viper.GetString("database.host")
	DB_PORT := viper.GetString("database.port")
	DB_DATABASE := viper.GetString("database.name")
	DB_USER := viper.GetString("database.username")
	DB_PASSWORD := viper.GetString("database.password")
	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE, DB_PORT)
	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})

	db.AutoMigrate(
		&auth.AccessToken{},
		&auth.DeviceToken{},
		&log.Log{},
		&user.User{},
	)

	if err != nil {
		Log.WithError(err).Error("failed to connect database")
	}

	connection, err := db.DB()
	if err != nil {
		Log.WithError(err).Error("failed to get connection")
	}
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(30)
	connection.SetConnMaxLifetime(time.Hour)
	return db
}
