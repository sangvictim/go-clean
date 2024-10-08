package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message     string `json:"message,omitempty"`
	Data        any    `json:"data,omitempty"`
	Errors      any    `json:"errors,omitempty"`
	Error       string `json:"error,omitempty"`
	CurrentPage int    `json:"currentPage,omitempty"`
	TotalPage   int    `json:"totalPage,omitempty"`
	Limit       int    `json:"limit,omitempty" validate:"min=1"`
}

type validationError struct {
	Namespace string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field     string `json:"field"`     // by passing alt name to ReportError like below
	Tag       string `json:"tag"`
	Message   string `json:"message"`
}

func ResponseJson(ctx echo.Context, code int, data Response) error {
	ctx.Response().Header().Set(echo.HeaderXContentTypeOptions, "nosniff")
	ctx.Response().Header().Set(echo.HeaderXFrameOptions, "deny")
	ctx.Response().Header().Set(echo.HeaderContentSecurityPolicy, "default-src 'none'")
	ctx.Response().Header().Set(echo.HeaderContentType, "application/json")
	ctx.Response().Header().Set(echo.HeaderVary, "Accept-Encoding")

	return ctx.JSON(code, data)
}

func HandleError(ctx echo.Context, err error) error {
	validationErrors := make([]validationError, 0)

	for _, err := range err.(validator.ValidationErrors) {
		e := validationError{
			Namespace: err.Namespace(),
			Field:     err.Field(),
			Tag:       err.Tag(),
			Message:   err.Error(),
		}
		validationErrors = append(validationErrors, e)

		_, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	return ctx.JSON(http.StatusBadRequest, Response{
		Message: "validation error",
		Errors:  validationErrors,
	})

}
