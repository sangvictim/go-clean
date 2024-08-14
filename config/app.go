package config

import (
	"go-clean/domain/auth"
	"go-clean/domain/log"
	"go-clean/domain/storage"
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
	Viper    *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup Repository
	authRepositorys := auth.NewAuthRepository(config.Log)
	userRepositorys := user.NewUserRepository(config.Log)
	storageController := storage.NewStorageController(config.Log, config.Viper)

	// setup Usecase
	authUsecase := auth.NewAuthUsecase(config.DB, config.Log, config.Validate, authRepositorys)
	userUsecase := user.NewUserUsecase(config.DB, config.Log, config.Validate, userRepositorys)

	// setup Controller
	authController := auth.NewAuthController(authUsecase, config.Log, config.Validate)
	userController := user.NewUserController(userUsecase, config.Log, config.Validate)

	// setup hook for logging to database
	config.Log.AddHook(&log.DBHook{DB: config.DB})

	// setup route
	routeConfig := routes.RouteConfig{
		App:               config.App,
		AuthController:    authController,
		UserController:    userController,
		StorageController: storageController,
	}

	routeConfig.Setup()
}
