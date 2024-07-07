package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zenoleg/shortener/internal/transport"
	"github.com/zenoleg/shortener/internal/transport/http/handler"
)

func NewEcho(shorten handler.ShortenHandler) *echo.Echo {
	e := echo.New()
	e.Binder = transport.NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.GET("/ping", func(ectx echo.Context) error {
		return ectx.String(http.StatusOK, "pong")
	})
	//e.GET("/link/:shortID", handler.Redirect)

	g := e.Group("/api/v1")
	g.POST("/shorten", shorten.Handle)

	//g.GET("/shorten", handler.GetShortURL)
	//g.GET("/original", handler.GetOriginal)

	return e
}
