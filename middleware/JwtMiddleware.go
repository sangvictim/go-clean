package middleware

import (
	"go-clean/pkg"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware(app *echo.Group) {
	app.Use(echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, pkg.Response{
				Message: "Unauthorized",
			})
		},
		SigningKey: []byte(os.Getenv("JWT_KEY")),
		ContextKey: "user",
	}))
}
