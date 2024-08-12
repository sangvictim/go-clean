package routes

import (
	"go-clean/domain/auth"
	"go-clean/domain/user"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App            *echo.Echo
	UserController *user.UserController
	AuthController *auth.AuthController
}

func (r *RouteConfig) Setup() {
	app := r.App.Group("/api")
	r.setUpGuest(app)
	r.setupAuth(app)
}

func (c *RouteConfig) setUpGuest(app *echo.Group) {
	// route for auth
	auth := app.Group("/auth")
	auth.POST("/register", c.AuthController.Register)
	auth.POST("/login", c.AuthController.Login)

	app.GET("/users", c.UserController.List)
	app.GET("/users/:id", c.UserController.Show)
	app.POST("/users", c.UserController.Create)
	app.PATCH("/users/:id", c.UserController.Update)
	app.DELETE("/users/:id", c.UserController.Delete)
}

func (c *RouteConfig) setupAuth(app *echo.Group) {

	app.GET("/user", func(c echo.Context) error {
		return c.JSON(200, "hello word")
	})
}
