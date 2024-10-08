package user

import (
	"go-clean/pkg"
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserService *UserService
	Log         *logrus.Logger
	Validate    *validator.Validate
}

func NewUserController(userService *UserService, log *logrus.Logger, validate *validator.Validate) *UserController {
	return &UserController{
		UserService: userService,
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

	response, total, err := c.UserService.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching user")
		return err
	}

	return pkg.ResponseJson(ctx, http.StatusOK, pkg.Response{
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

	response, err := c.UserService.FindById(ctx.Request().Context(), id)
	if err != nil {
		c.Log.WithError(err).Error("error getting user")
		return err
	}

	if response.Id == 0 {
		return pkg.ResponseJson(ctx, http.StatusNotFound, pkg.Response{
			Message: "user not found",
		})
	}

	return pkg.ResponseJson(ctx, http.StatusOK, pkg.Response{
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
		pkg.HandleError(ctx, err)
	}

	request := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   &user.Avatar,
	}

	response, err := c.UserService.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error(err)
		return err
	}

	return pkg.ResponseJson(ctx, http.StatusCreated, pkg.Response{
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
		pkg.HandleError(ctx, err)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	request := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   &user.Avatar,
	}

	response, err := c.UserService.Update(ctx.Request().Context(), request, id)
	if err != nil {
		c.Log.WithError(err).Error(err.Error())
		return err
	}

	return pkg.ResponseJson(ctx, http.StatusOK, pkg.Response{
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

	if err := c.UserService.Delete(ctx.Request().Context(), id); err != nil {
		c.Log.WithError(err).Error(err.Error())
		return err
	}

	return pkg.ResponseJson(ctx, http.StatusOK, pkg.Response{
		Message: "user deleted",
	})
}

// @tags			User
// @summary		Current User
// @description	Current User
// @Accept			json
// @Produce		json
// @Success		200	{object} UserDetail "Current User"
// @Router			/users/profile [get]		true	"Current User"
// @Security Bearer
func (c *UserController) Profile(ctx echo.Context) error {
	userId := c.currentUser(ctx)["id"].(float64)

	user, err := c.UserService.FindById(ctx.Request().Context(), int(userId))
	if err != nil {
		return err
	}
	return pkg.ResponseJson(ctx, http.StatusOK, pkg.Response{
		Message: "Current User",
		Data:    user,
	})
}

func (c *UserController) currentUser(ctx echo.Context) jwt.MapClaims {
	user := ctx.Get("user").(*jwt.Token)
	result := user.Claims.(jwt.MapClaims)

	return result
}
