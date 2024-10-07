package config

import (
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewSwaggerConfig(app *echo.Echo) {
	if os.Getenv("APP_ENV") == "local" {
		app.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
