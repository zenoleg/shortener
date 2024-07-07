package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
	http2 "github.com/zenoleg/shortener/internal/transport/http"
	"github.com/zenoleg/shortener/internal/transport/http/handler/mocks"
)

func TestRedirectHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When original not found then return 404", func(t *testing.T) {
		redirectUseCase := mocks.NewGetOriginalForRedirect(t)

		redirectUseCase.
			On("Do", mock.Anything).
			Return(domain.URL(""), storage.ErrURLNotFound).
			Once()

		handler := NewRedirectHandler(redirectUseCase, zerolog.Logger{})

		e := makeRedirectTestEnv(t, handler)

		e.GET("/link/123").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("When found then return 301", func(t *testing.T) {
		redirectUseCase := mocks.NewGetOriginalForRedirect(t)

		redirectUseCase.
			On("Do", mock.Anything).
			Return(domain.URL("https://example.com"), nil).
			Once()

		handler := NewRedirectHandler(redirectUseCase, zerolog.Logger{})

		e := makeRedirectTestEnv(t, handler)

		e.GET("/link/123").
			Expect().
			Status(http.StatusMovedPermanently).
			Header("Location").
			IsEqual("https://example.com")
	})
}

func makeRedirectTestEnv(t *testing.T, handler RedirectHandler) *httpexpect.Expect {
	e := echo.New()
	e.Binder = http2.NewValidationBinder()
	e.HideBanner = true

	e.Use(middleware.Recover())

	e.GET("/link/:id", handler.Handle)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(e),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
