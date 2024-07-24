package config

import (
	"go-clean/internal/controller"
	"go-clean/internal/repository"
	"go-clean/internal/usecase"
	"go-clean/routes"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup Repository
	userRepository := repository.NewUserRepository(config.Log)

	// setup Usecase
	userUsecase := usecase.NewUserUsecase(config.DB, config.Log, config.Validate, userRepository)

	// setup Controller
	userController := controller.NewUserController(userUsecase, config.Log)

	// setup route
	routeConfig := routes.RouteConfig{
		App:            config.App,
		UserController: userController,
	}

	routeConfig.Setup()
}
