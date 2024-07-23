package server

import (
	"go-clean/routes"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	echo *echo.Echo
	db   *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		echo: echo.New(),
		db:   db,
	}
}

func (s *Server) Start() error {

	// setup route
	r := routes.NewRoute(s.echo, s.db)
	r.Setup()

	// start server
	s.echo.Logger.Fatal(s.echo.Start(":8080"))
	return nil
}
