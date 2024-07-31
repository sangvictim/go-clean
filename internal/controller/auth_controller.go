package controller

import (
	"go-clean/internal/model"
	"go-clean/internal/usecase"
	apiResponse "go-clean/utils/api_response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthUsecase *usecase.AuthUsecase
	Log         *logrus.Logger
}

func NewAuthController(usecase *usecase.AuthUsecase, log *logrus.Logger) *AuthController {
	return &AuthController{
		AuthUsecase: usecase,
		Log:         log,
	}
}

// @tags			Auth
// @summary		Register User
// @description	Register User
// @Accept			json
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router			/auth/register [post]
// @Param request body model.UserRequest true "register user"
func (c *AuthController) Register(ctx echo.Context) error {
	user := new(model.User)

	if err := ctx.Bind(user); err != nil {
		c.Log.WithError(err).Error("error validating user")
		return ctx.JSON(http.StatusBadRequest, apiResponse.Response{
			Message: "error validating user",
			Errors:  err.Error(),
		})
	}

	request := &model.UserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	response, _ := c.AuthUsecase.Register(ctx.Request().Context(), request)

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "register success",
		Data:    response,
	})
}

func (c *AuthController) Login(ctx echo.Context) error {
	user := new(model.LoginRequest)

	if err := ctx.Bind(user); err != nil {
		c.Log.WithError(err).Error("error validating login")
		return ctx.JSON(http.StatusBadRequest, apiResponse.Response{
			Message: "error validating login",
			Errors:  err.Error(),
		})
	}

	request := &model.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	response, _ := c.AuthUsecase.Login(ctx.Request().Context(), request)

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "login success",
		Data:    response,
	})
}
