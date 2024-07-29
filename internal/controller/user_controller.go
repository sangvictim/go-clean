package controller

import (
	"go-clean/internal/model"
	"go-clean/internal/usecase"
	apiResponse "go-clean/pkg/api_response"

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

// @tags			User
// @summary		List User
// @description	List User
// @Accept			json
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router			/users [get]
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
	return ctx.JSON(200, apiResponse.Response{
		Data:    response,
		Message: "success",
	})
}
