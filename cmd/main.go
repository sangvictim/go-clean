package main

import (
	"flag"
	"go-clean/config"
	seeder "go-clean/seed"
	"log"
	"os"

	_ "go-clean/docs"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := echo.New()
	validate := validator.New()
	log := config.NewLogger()
	db := config.NewDatabase(log)
	config.NewSwaggerConfig(app)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})

	// seed
	seedFlag := flag.Bool("seed", false, "seed database")
	flag.Parse()
	if *seedFlag {
		seeder.DatabaseSeeder(db)
	}
	app.Logger.Fatal(app.Start(os.Getenv("APP_URL")))
}
