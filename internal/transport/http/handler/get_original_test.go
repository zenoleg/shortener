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
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
	"github.com/zenoleg/shortener/internal/transport/http/handler/mocks"
)

func TestGetOriginalURLHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When request is invalid, then return 400", func(t *testing.T) {
		getOriginalURL := mocks.NewGetOriginalUseCase(t)
		getOriginalURL.AssertNumberOfCalls(t, "Do", 0)

		handler := NewGetOriginalURLHandler(getOriginalURL, zerolog.Logger{})

		e := makeGetOriginalURLTestEnv(t, handler)

		e.GET("/api/v1/shorten").
			WithQuery("url", "").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("When use case fails, then return 500", func(t *testing.T) {
		getOriginalURL := mocks.NewGetOriginalUseCase(t)
		getOriginalURL.
			On("Do", mock.Anything, mock.Anything).
			Return(domain.URL(""), assert.AnError)

		handler := NewGetOriginalURLHandler(getOriginalURL, zerolog.Logger{})

		e := makeGetOriginalURLTestEnv(t, handler)

		e.GET("/api/v1/shorten").
			WithQuery("url", "https://example.com").
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("When url not found, then return 404", func(t *testing.T) {
		getOriginalURL := mocks.NewGetOriginalUseCase(t)
		getOriginalURL.
			On("Do", mock.Anything, mock.Anything).
			Return(domain.URL(""), storage.ErrURLNotFound).
			Once()

		handler := NewGetOriginalURLHandler(getOriginalURL, zerolog.Logger{})

		e := makeGetOriginalURLTestEnv(t, handler)

		e.GET("/api/v1/shorten").
			WithQuery("url", "https://example.com").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("When use case succeeds, then return 200", func(t *testing.T) {
		getOriginalURL := mocks.NewGetOriginalUseCase(t)
		getOriginalURL.
			On("Do", mock.Anything, mock.Anything).
			Return(domain.URL("https://example.com"), nil)

		handler := NewGetOriginalURLHandler(getOriginalURL, zerolog.Logger{})

		e := makeGetOriginalURLTestEnv(t, handler)

		e.GET("/api/v1/shorten").
			WithQuery("url", "https://local.com/link/123").
			Expect().
			Status(http.StatusOK).
			JSON().Object().HasValue("destination", "https://example.com")
	})
}
func makeGetOriginalURLTestEnv(t *testing.T, handler GetOriginalURLHandler) *httpexpect.Expect {
	e := echo.New()
	e.Binder = NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.GET("/api/v1/shorten", handler.Handle)

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
