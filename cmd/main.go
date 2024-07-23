package main

import (
	"go-clean/pkg/database"
	"go-clean/server"

	"github.com/labstack/gommon/log"
)

func main() {

	db, _ := database.NewDatabaseConnection()

	s := server.NewServer(db)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
