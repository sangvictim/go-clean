package main

import (
	"fmt"
	"go-clean/config"

	_ "go-clean/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title			Go-Clean Example API
// @version		1.0
// @description	This is a sample api.
// @host		localhost:8080
// @BasePath	/api
func main() {
	viperConfig := config.NewViper()
	app := echo.New()
	validate := validator.New()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", viperConfig.GetInt("api.port"))))
}
