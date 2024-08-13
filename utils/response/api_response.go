package apiResponse

import "github.com/labstack/echo/v4"

type Response struct {
	Message     string `json:"message,omitempty"`
	Data        any    `json:"data,omitempty"`
	Errors      any    `json:"errors,omitempty"`
	Error       string `json:"error,omitempty"`
	CurrentPage int    `json:"currentPage,omitempty"`
	TotalPage   int    `json:"totalPage,omitempty"`
	Limit       int    `json:"limit,omitempty" validate:"min=1"`
}

func ResponseJson(ctx echo.Context, code int, data Response) error {
	ctx.Response().Header().Set(echo.HeaderXContentTypeOptions, "nosniff")
	ctx.Response().Header().Set(echo.HeaderXFrameOptions, "deny")
	ctx.Response().Header().Set(echo.HeaderContentSecurityPolicy, "default-src 'none'")
	ctx.Response().Header().Set(echo.HeaderContentType, "application/json")
	ctx.Response().Header().Set(echo.HeaderVary, "Accept-Encoding")

	return ctx.JSON(code, data)
}
