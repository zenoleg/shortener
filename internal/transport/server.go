package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())

	return e
}
