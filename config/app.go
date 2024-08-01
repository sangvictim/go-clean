package config

import (
	"go-clean/internal/controller"
	"go-clean/internal/repository"
	"go-clean/internal/usecase"
	"go-clean/routes"

	"github.com/go-playground/validator/v10"
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
	authUsecase := usecase.NewAuthUsecase(config.DB, config.Log, config.Validate, userRepository)
	userUsecase := usecase.NewUserUsecase(config.DB, config.Log, config.Validate, userRepository)

	// setup Controller
	authController := controller.NewAuthController(authUsecase, config.Log)
	userController := controller.NewUserController(userUsecase, config.Log)

	// setup hook for logging to database
	config.Log.AddHook(&repository.DBHook{DB: config.DB})

	// setup route
	routeConfig := routes.RouteConfig{
		App:            config.App,
		AuthController: authController,
		UserController: userController,
	}

	routeConfig.Setup()
}
