package database

import (
	"go-clean/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
	DB_DATABASE = "golang"
	DB_USER     = "postgres"
	DB_PASSWORD = "123456"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	var conn = "host=" + DB_HOST + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_DATABASE + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	db.AutoMigrate(&models.User{})
	connection, _ := db.DB()

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)
	return db, err
}
