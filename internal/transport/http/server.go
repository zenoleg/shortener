package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zenoleg/shortener/internal/transport/http/handler"
)

func NewEcho(
	shorten handler.ShortenHandler,
	getShorten handler.GetShortURLHandler,
	getOriginal handler.GetOriginalURLHandler,
	redirect handler.RedirectHandler,
) *echo.Echo {
	e := echo.New()
	e.Binder = handler.NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.GET("/ping", func(ectx echo.Context) error {
		return ectx.String(http.StatusOK, "pong")
	})
	e.GET("/link/:shortID", redirect.Handle)

	g := e.Group("/api/v1")
	g.POST("/shorten", shorten.Handle)

	g.GET("/shorten", getShorten.Handle)
	g.GET("/original", getOriginal.Handle)

	return e
}
