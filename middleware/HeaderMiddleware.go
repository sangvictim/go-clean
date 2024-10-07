package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// TODO: nned to fix it
func HeaderMiddleware(app *echo.Echo) {
	// app.Use(middleware.Logger())
	app.Use(middleware.Secure())
	app.Use(middleware.Recover())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.PATCH, echo.POST, echo.DELETE},
	}))
	if os.Getenv("APP_ENV") == "dev" {
		app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
			Skipper: func(c echo.Context) bool {
				if strings.Contains(c.Request().URL.Path, "swagger") {
					return true
				}
				return false
			},
		}))
	}

	// rate limiter
	configRateLimiter := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 20, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	app.Use(middleware.RateLimiterWithConfig(configRateLimiter))

	//time out
	app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Connection timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 30 * time.Second,
	}))

}
