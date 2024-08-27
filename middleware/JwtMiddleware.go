package middleware

import (
	apiResponse "go-clean/utils/response"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func JwtMiddleware(app *echo.Group, viper *viper.Viper) {
	app.Use(echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, apiResponse.Response{
				Message: "Unauthorized",
			})
		},
		SigningKey: []byte(viper.GetString("jwt.key")),
		ContextKey: "user",
	}))
}
