package config

import (
	"fmt"
	"go-clean/domain/log"
	"go-clean/domain/user"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(viper *viper.Viper, Log *logrus.Logger) *gorm.DB {
	DB_DRIVER := viper.GetString("database.driver")
	DB_HOST := viper.GetString("database.host")
	DB_PORT := viper.GetString("database.port")
	DB_DATABASE := viper.GetString("database.name")
	DB_USER := viper.GetString("database.username")
	DB_PASSWORD := viper.GetString("database.password")
	var db *gorm.DB
	var err error

	if DB_DRIVER == "postgres" {
		conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE, DB_PORT)
		db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_DATABASE)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	db.AutoMigrate(
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
	connection.SetMaxOpenConns(10)
	return db
}
