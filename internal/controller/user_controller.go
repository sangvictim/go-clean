package controller

import (
	"go-clean/internal/model"
	"go-clean/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUsecase *usecase.UserUsecase
	Log         *logrus.Logger
}

func NewUserController(userUsecase *usecase.UserUsecase, log *logrus.Logger) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
		Log:         log,
	}
}

func (c *UserController) List(ctx echo.Context) error {

	request := &model.UserSearchRequest{
		Name:  ctx.QueryParam("name"),
		Email: ctx.QueryParam("email"),
	}

	response, err := c.UserUsecase.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return err
	}
	return ctx.JSON(200, response)
}
