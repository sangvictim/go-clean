package config

import (
	repositoryLog "go-clean/domain/log/repository"
	userController "go-clean/domain/user/controller"
	userRepository "go-clean/domain/user/repository"
	userUsecase "go-clean/domain/user/usecase"
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
	userRepositorys := userRepository.NewUserRepository(config.Log)

	// setup Usecase
	// authUsecase := usecase.NewAuthUsecase(config.DB, config.Log, config.Validate, userRepository)
	userUsecase := userUsecase.NewUserUsecase(config.DB, config.Log, config.Validate, userRepositorys)

	// setup Controller
	// authController := controller.NewAuthController(authUsecase, config.Log)
	userController := userController.NewUserController(userUsecase, config.Log)

	// setup hook for logging to database
	config.Log.AddHook(&repositoryLog.DBHook{DB: config.DB})

	// setup route
	routeConfig := routes.RouteConfig{
		App: config.App,
		// AuthController: authController,
		UserController: userController,
	}

	routeConfig.Setup()
}
