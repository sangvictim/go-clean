package main

import (
	"flag"
	"fmt"
	"go-clean/config"
	"go-clean/middleware"
	seeder "go-clean/seed"

	_ "go-clean/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// @title			Go-Clean Example API
// @version		1.0
// @description	This is a sample api.
// @host		localhost:8080
// @BasePath	/api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	viperConfig := config.NewViper()
	app := echo.New()
	validate := validator.New()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	config.NewSwaggerConfig(app, viperConfig)
	middleware.HeaderMiddleware(app)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Viper:    viperConfig,
	})

	// seed
	seedFlag := flag.String("seed", "true", "seed database")
	if *seedFlag == "true" {
		seeder.DatabaseSeeder(db)
	}

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", viperConfig.GetInt("api.port"))))
}
