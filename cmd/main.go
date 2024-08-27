package main

import (
	"flag"
	"fmt"
	"go-clean/config"
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
// @description Enter the token with the `Bearer` prefix, e.g. "Bearer abcde12345"
func main() {
	viperConfig := config.NewViper()
	app := echo.New()
	validate := validator.New()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	config.NewSwaggerConfig(app, viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Viper:    viperConfig,
	})

	// seed
	seedFlag := flag.Bool("seed", false, "seed database")
	flag.Parse()
	if *seedFlag {
		seeder.DatabaseSeeder(db)
	}

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", viperConfig.GetInt("api.port"))))
}
