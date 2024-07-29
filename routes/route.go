package routes

import (
	"go-clean/internal/controller"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App            *echo.Echo
	UserController *controller.UserController
}

func (r *RouteConfig) Setup() {
	app := r.App.Group("/api")
	r.setUpGuest(app)
	r.setupAuth(app)
}

func (c *RouteConfig) setUpGuest(app *echo.Group) {

	app.GET("/users", c.UserController.List)
	app.GET("/users/:id", c.UserController.Show)
	app.POST("/users", c.UserController.Create)
}

func (c *RouteConfig) setupAuth(app *echo.Group) {

	app.GET("/user", func(c echo.Context) error {
		return c.JSON(200, "hello word")
	})
}
