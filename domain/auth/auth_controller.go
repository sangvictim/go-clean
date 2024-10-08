package auth

import (
	"go-clean/pkg"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthService *AuthService
	Log         *logrus.Logger
	Validate    *validator.Validate
}

func NewAuthController(authService *AuthService, log *logrus.Logger, validate *validator.Validate) *AuthController {
	return &AuthController{
		AuthService: authService,
		Log:         log,
		Validate:    validate,
	}
}

// @tags			Auth
// @summary		Register User
// @description	Register User
// @Accept			json
// @Produce		json
// @Success		200		{string}	string
// @Router			/auth/register [post]
// @Param request body Register true "user register"
func (c *AuthController) Register(ctx echo.Context) error {
	user := new(Register)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		return pkg.HandleError(ctx, err)
	}

	request := &Register{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := c.AuthService.Register(ctx.Request().Context(), request); err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, pkg.Response{
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
// @Param request body LoginRequest true "user login"
func (c *AuthController) Login(ctx echo.Context) error {
	user := new(LoginRequest)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		return pkg.HandleError(ctx, err)
	}

	request := &LoginRequest{
		Email:      user.Email,
		Password:   user.Password,
		DeviceID:   ctx.Request().Header.Get("X-Device-Id"),
		DeviceType: ctx.Request().Header.Get("X-Device-Type"),
		UserAgent:  ctx.Request().UserAgent(),
	}

	response, err := c.AuthService.Login(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return ctx.JSON(http.StatusOK, pkg.Response{
		Message: "Login success",
		Data:    response,
	})
}

func (c *AuthController) Logout(ctx echo.Context) error {
	refresh_token := new(LogoutRequest)
	if err := ctx.Bind(refresh_token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(refresh_token); err != nil {
		return pkg.HandleError(ctx, err)
	}

	req := &AccessToken{
		RefreshToken: refresh_token.RefreshToken,
	}
	getDevice := ctx.Request().Header.Get("X-Device-Id")

	if err := c.AuthService.Logout(ctx.Request().Context(), req.RefreshToken, getDevice); err != nil {
		return ctx.JSON(http.StatusUnauthorized, pkg.Response{
			Message: "Logout failed",
		})
	}

	return pkg.ResponseJson(ctx, http.StatusOK,
		pkg.Response{
			Message: "Logout success",
		},
	)
}
