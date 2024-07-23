package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Route struct {
	echo *echo.Echo
	db   *gorm.DB
}

func NewRoute(e *echo.Echo, db *gorm.DB) *Route {
	return &Route{
		echo: e,
		db:   db,
	}
}

func (r *Route) Setup() {
	app := r.echo.Group("/api")
	setUpGuest(app)
	setupAuth(app)
}

func setUpGuest(app *echo.Group) {

	app.GET("/users", func(c echo.Context) error {
		return c.JSON(200, "hello word")
	})
}

func setupAuth(app *echo.Group) {

	app.GET("/user", func(c echo.Context) error {
		return c.JSON(200, "hello word")
	})
}
