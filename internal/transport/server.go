package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho(handler ShortenHandler) *echo.Echo {
	e := echo.New()
	e.Binder = NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.GET("/ping", func(ectx echo.Context) error {
		return ectx.String(http.StatusOK, "pong")
	})

	g := e.Group("/api/v1")
	g.GET("/shorten", handler.GenerateShortenURL)
	g.POST("/shorten", handler.Shorten)
	g.GET("/original", handler.GetOriginal)

	return e
}
