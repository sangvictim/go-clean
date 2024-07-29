package controller

import (
	"go-clean/internal/model"
	"go-clean/internal/usecase"
	apiResponse "go-clean/utils/api_response"
	"math"
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
