package user

import (
	FormValidator "go-clean/utils/formValidate"
	apiResponse "go-clean/utils/response"
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUsecase *UserUsecase
	Log         *logrus.Logger
	Validate    *validator.Validate
}

func NewUserController(userUsecase *UserUsecase, log *logrus.Logger, validate *validator.Validate) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
		Log:         log,
		Validate:    validate,
	}
}

// @tags			User
// @summary		List User
// @description	List User
// @Accept			json
// @Produce		json
// @Success		200		{object}	UserDetail
// @Router			/users [get]
// @param			search		query		string		false	"serch email or name"
// @param			page		query		int		false	"page"
// @param			size		query		int		false	"size"
// @param			orderBy		query	string false "orderBy"
// @param			orderDirection		query string false "orderDirection"
// @Security Bearer
func (c *UserController) List(ctx echo.Context) error {

	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

	request := &UserSearchRequest{
		Search:         ctx.QueryParam("search"),
		Page:           page,
		Limit:          limit,
		OrderBy:        ctx.QueryParam("orderBy"),
		OrderDirection: ctx.QueryParam("orderDirection"),
	}

	if request.Limit == 0 {
		request.Limit = 10
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.OrderBy == "" {
		request.OrderBy = "created_at"
	}
	if request.OrderDirection == "" {
		request.OrderDirection = "desc"
	}

	response, total, err := c.UserUsecase.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching user")
		return err
	}

	return ctx.JSON(200, apiResponse.Response{
		Data:        response,
		Message:     "List User",
		CurrentPage: request.Page,
		TotalPage:   int(math.Ceil(float64(total) / float64(request.Limit))),
		Limit:       request.Limit,
	})
}

// @tags			User
// @summary		Detail User
// @description	Detail User
// @Accept			json
// @Produce		json
// @Success		200		{object}	UserDetail
// @Router			/users/{id} [get]
// @param			id		path		string		true	"User ID"
// @Security Bearer
func (c *UserController) Show(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	response, err := c.UserUsecase.FindById(ctx.Request().Context(), id)
	if err != nil {
		c.Log.WithError(err).Error("error getting user")
		return err
	}

	if response.Id == 0 {
		return ctx.JSON(http.StatusNotFound, apiResponse.Response{
			Message: "user not found",
		})
	}

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "Detail User",
		Data:    response,
	})
}

// @tags			User
// @summary		Create User
// @description	Create User
// @Accept			json
// @Produce		json
// @Success		200		{object}	UserDetail
// @Router			/users [post]
// @Security Bearer
// @Param request body UserCreate true "user request"
func (c *UserController) Create(ctx echo.Context) error {
	user := new(UserCreate)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		FormValidator.HandleError(ctx, err)
	}

	request := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	response, err := c.UserUsecase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, apiResponse.Response{
		Message: "success",
		Data:    response,
	})
}

// @tags			User
// @summary		Update User
// @description	Update User
// @Accept			json
// @Produce		json
// @Success		200		{object}	UserDetail
// @Router			/users/{id} [PATCH]
// @param			id		path		string		true	"User ID"
// @Param request body UserUpdate true "user request"
// @Security Bearer
func (c *UserController) Update(ctx echo.Context) error {
	user := new(UserUpdate)

	if err := ctx.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate.Struct(user); err != nil {
		FormValidator.HandleError(ctx, err)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	request := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	response, err := c.UserUsecase.Update(ctx.Request().Context(), request, id)
	if err != nil {
		c.Log.WithError(err).Error(err.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse.Response{
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
// @Security Bearer
func (c *UserController) Delete(ctx echo.Context) error {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := c.UserUsecase.Delete(ctx.Request().Context(), id); err != nil {
		c.Log.WithError(err).Error(err.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse.Response{
		Message: "user deleted",
	})
}
