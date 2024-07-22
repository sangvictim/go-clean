package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_CONNECTION = "pgsql" // pgsql or mysql
	DB_HOST       = "localhost"
	DB_PORT       = "5432"
	DB_DATABASE   = "golang"
	DB_USER       = "postgres"
	DB_PASSWORD   = "123456"
)

func NewDatabaseConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(database(DB_CONNECTION)), &gorm.Config{})

	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		panic(err)
	}

	db.AutoMigrate()
	connection, _ := db.DB()

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)
	return db
}

func database(db_connection string) string {
	var conn string
	switch db_connection {

	case "pgsql":
		conn = "host=" + DB_HOST + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_DATABASE + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Jakarta"

	case "mysql":
		conn = DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return conn
}
