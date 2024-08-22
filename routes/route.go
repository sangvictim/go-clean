package routes

import (
	"go-clean/domain/auth"
	"go-clean/domain/storage"
	"go-clean/domain/user"
	"go-clean/middleware"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App               *echo.Echo
	Viper             *viper.Viper
	UserController    *user.UserController
	AuthController    *auth.AuthController
	StorageController *storage.UploadController
}

func (r *RouteConfig) Setup() {
	app := r.App.Group("/api")
	middleware.HeaderMiddleware(r.App)
	r.setUpGuest(app)
	r.setupAuth(app)
}

func (c *RouteConfig) setUpGuest(app *echo.Group) {
	// route for auth
	guest := app.Group("/auth")
	guest.POST("/register", c.AuthController.Register)
	guest.POST("/login", c.AuthController.Login)

	// route for public
	app.GET("/public/:key", c.StorageController.GetFile)
}

func (c *RouteConfig) setupAuth(app *echo.Group) {

	// Middleware for auth with jwt
	middleware.JwtMiddleware(app, c.Viper)

	app.GET("/users", c.UserController.List)
	app.GET("/users/:id", c.UserController.Show)
	app.POST("/users", c.UserController.Create)
	app.PATCH("/users/:id", c.UserController.Update)
	app.DELETE("/users/:id", c.UserController.Delete)
	app.POST("/upload", c.StorageController.UploadFile)

	app.POST("/logout", c.AuthController.Logout)

}
