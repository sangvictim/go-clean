package config

import (
	"go-clean/domain/log"
	"go-clean/domain/user"
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

// TODO: fix it
func Bootstrap(config *BootstrapConfig) {
	// setup Repository
	userRepositorys := user.NewUserRepository(config.Log)

	// setup Usecase
	// authUsecase := usecase.NewAuthUsecase(config.DB, config.Log, config.Validate, userRepository)
	userUsecase := user.NewUserUsecase(config.DB, config.Log, config.Validate, userRepositorys)

	// setup Controller
	// authController := controller.NewAuthController(authUsecase, config.Log)
	userController := user.NewUserController(userUsecase, config.Log, config.Validate)

	// setup hook for logging to database
	config.Log.AddHook(&log.DBHook{DB: config.DB})

	// setup route
	routeConfig := routes.RouteConfig{
		App: config.App,
		// AuthController: authController,
		UserController: userController,
	}

	routeConfig.Setup()
}
