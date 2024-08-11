package FormValidator

import (
	"encoding/json"
	"fmt"
	apiResponse "go-clean/utils/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type validationError struct {
	Namespace string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field     string `json:"field"`     // by passing alt name to ReportError like below
	Tag       string `json:"tag"`
	Message   string `json:"message"`
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
	return ctx.JSON(http.StatusBadRequest, apiResponse.Response{
		Message: "validation error",
		Errors:  validationErrors,
	})

}
