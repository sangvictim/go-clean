package main

import (
	"fmt"
	"go-clean/config"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

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

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", viperConfig.GetInt("api.port"))))
}
