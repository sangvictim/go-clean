package controller

import (
	"go-clean/internal/model"
	"go-clean/internal/usecase"
	apiResponse "go-clean/utils/api_response"
	"math"
	"net/http"
	"strconv"

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
// @param			page		query		int		false	"page"
// @param			size		query		int		false	"size"
func (c *UserController) List(ctx echo.Context) error {

	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	size, _ := strconv.Atoi(ctx.QueryParam("size"))

	request := &model.UserSearchRequest{
		Name:  ctx.QueryParam("name"),
		Email: ctx.QueryParam("email"),
		Page:  page,
		Size:  size,
	}

	response, total, err := c.UserUsecase.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching user")
		return err
	}

	return ctx.JSON(200, apiResponse.Response{
		Data:        response,
		Message:     "success",
		TotalPage:   int(math.Ceil(float64(total) / float64(size))),
		Size:        request.Size,
		CurrentPage: request.Page,
	})
}

// @tags			User
// @summary		Detail User
// @description	Detail User
// @Accept			json
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router			/users/{id} [get]
// @param			id		path		string		true	"User ID"
func (c *UserController) Show(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	request := &model.UserId{
		ID: id,
	}

	response, err := c.UserUsecase.FindById(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting user")
		return err
	}

	return ctx.JSON(200, apiResponse.Response{
		Message: "Detail User",
		Data:    response,
	})
}

// @tags			User
// @summary		Create User
// @description	Create User
// @Accept			json
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router			/users [post]
// @Param request body model.UserRequest true "user request"
func (c *UserController) Create(ctx echo.Context) error {
	user := new(model.User)

	if err := ctx.Bind(user); err != nil {
		c.Log.WithError(err).Error("error binding user")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	request := &model.UserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	response, err := c.UserUsecase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, apiResponse.Response{
		Message: "user created",
		Data:    response,
	})
}

// @tags			User
// @summary		Update User
// @description	Update User
// @Accept			json
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router			/users/{id} [put]
// @param			id		path		string		true	"User ID"
// @Param request body model.UserRequest true "user request"
func (c *UserController) Update(ctx echo.Context) error {
	user := new(model.User)

	if err := ctx.Bind(user); err != nil {
		c.Log.WithError(err).Error("error binding user")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	request := &model.UserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	response, err := c.UserUsecase.Update(ctx.Request().Context(), request, id)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, apiResponse.Response{
		Message: "user updated",
		Data:    response,
	})
}

// @tags			User
// @summary		Delete User
// @description	Delete User
// @Accept			json
// @Produce		json
// @Success		200		{string}	string	"user deleted"
// @Router			/users/{id} [delete]
// @param			id		path		string		true	"User ID"
func (c *UserController) Delete(ctx echo.Context) error {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := c.UserUsecase.Delete(ctx.Request().Context(), id); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "user deleted",
	})
}
