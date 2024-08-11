package config

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewSwaggerConfig(app *echo.Echo, config *viper.Viper) {
	if config.GetString("api.deploy") == "dev" {
		app.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
