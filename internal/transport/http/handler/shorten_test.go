package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zenoleg/shortener/internal/transport/http/handler/mocks"
	"github.com/zenoleg/shortener/internal/usecase"
)

func TestShortenHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When request is invalid, then return 400", func(t *testing.T) {
		shortenUseCase := mocks.NewShortenUseCase(t)
		shortenUseCase.AssertNumberOfCalls(t, "Do", 0)

		handler := NewShortenHandler(shortenUseCase, zerolog.Logger{})

		e := makeShortenTestEnv(t, handler)

		e.POST("/api/v1/shorten").
			WithJSON(ShortenRequest{URL: ""}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("When use case fails, then return 500", func(t *testing.T) {
		shortenUseCase := mocks.NewShortenUseCase(t)
		shortenUseCase.
			On("Do", mock.Anything, mock.Anything).
			Return(usecase.DestinationURL(""), assert.AnError)

		handler := NewShortenHandler(shortenUseCase, zerolog.Logger{})

		e := makeShortenTestEnv(t, handler)

		e.POST("/api/v1/shorten").
			WithJSON(ShortenRequest{URL: "https://example.com"}).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("When use case succeeds, then return 200", func(t *testing.T) {
		shortenUseCase := mocks.NewShortenUseCase(t)
		shortenUseCase.
			On("Do", mock.Anything, mock.Anything).
			Return(usecase.DestinationURL("https://example.com"), nil)

		handler := NewShortenHandler(shortenUseCase, zerolog.Logger{})

		e := makeShortenTestEnv(t, handler)

		e.POST("/api/v1/shorten").
			WithJSON(ShortenRequest{URL: "https://example.com"}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().HasValue("destination", "https://example.com")
	})
}

func makeShortenTestEnv(t *testing.T, handler ShortenHandler) *httpexpect.Expect {
	e := echo.New()
	e.Binder = NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.POST("/api/v1/shorten", handler.Handle)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(e),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
