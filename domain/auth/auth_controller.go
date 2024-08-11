package auth

import (
	FormValidator "go-clean/utils/formValidate"
	apiResponse "go-clean/utils/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthUsecase *AuthUsecase
	Log         *logrus.Logger
	Validate    *validator.Validate
}

func NewAuthController(authUsecase *AuthUsecase, log *logrus.Logger, validate *validator.Validate) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
		Log:         log,
		Validate:    validate,
	}
}

// @tags			Auth
// @summary		Register User
// @description	Register User
// @Accept			json
// @Produce		json
// @Success		200		{object}	Register
// @Router			/auth/register [post]
// @Param request body Register true "user register"

// TODO: fix return validation error
func (c *AuthController) Register(ctx echo.Context) error {
	user := new(Register)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		FormValidator.HandleError(ctx, err)
	}

	request := &Register{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	_, err := c.AuthUsecase.Register(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, apiResponse.Response{
		Message: "Register success",
	})
}

// @tags			Auth
// @summary		User Login
// @description	User Login
// @Accept			json
// @Produce		json
// @Success		200		{object}	LoginResponse
// @Router			/auth/login [post]
// @Param request body LoginResponse true "user login"
func (c *AuthController) Login(ctx echo.Context) error {
	user := new(LoginRequest)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		FormValidator.HandleError(ctx, err)
	}

	request := &LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	response, err := c.AuthUsecase.Login(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "Login success",
		Data:    response,
	})
}
