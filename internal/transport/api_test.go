package transport

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/shortener"
)

func TestShortenHandler_Shorten(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(initHandler(map[string]string{})),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testShorten(e)
}

func TestShortenHandler_GetShortURL(t *testing.T) {
	store := map[string]string{
		"t92YuUGbw1WY4V2LvozcwRHdoB": "https://example.com",
	}

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(initHandler(store)),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testGetShortURL(e)
}

func TestShortenHandler_GetOriginal(t *testing.T) {
	store := map[string]string{
		"t92YuUGbw1WY4V2LvozcwRHdoB": "https://example.com",
	}

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(initHandler(store)),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testGetOriginal(e)
}

func TestShortenHandler_Redirect(t *testing.T) {
	store := map[string]string{
		"t92YuUGbw1WY4V2LvozcwRHdoB": "https://example.com",
	}

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(initHandler(store)),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testRedirect(e)
}

func testShorten(e *httpexpect.Expect) {
	e.POST("/api/v1/shorten").
		WithJSON(ShortenRequest{URL: "invalid-url"}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST("/api/v1/shorten").
		WithJSON(ShortenRequest{URL: "https://example.com"}).
		WithHost("service.com").
		Expect().
		Status(http.StatusCreated).
		JSON().Object().HasValue("destination", "http://service.com/link/t92YuUGbw1WY4V2LvozcwRHdoB")
}

func testGetShortURL(e *httpexpect.Expect) {
	// When URL did not pass validation, then return 400
	e.GET("/api/v1/shorten").
		WithQuery("url", "invalid-url").
		Expect().
		Status(http.StatusBadRequest)

	// When valid URL passed but not found, then return 404
	e.GET("/api/v1/shorten").
		WithQuery("url", "https://example.ru").
		Expect().
		Status(http.StatusNotFound)

	// When valid URL passed and found, then return 200
	e.GET("/api/v1/shorten").
		WithQuery("url", "https://example.com").
		WithHost("service.com").
		Expect().
		Status(http.StatusOK).
		JSON().Object().HasValue("destination", "http://service.com/link/t92YuUGbw1WY4V2LvozcwRHdoB")
}

func testGetOriginal(e *httpexpect.Expect) {
	// When URL did not pass validation, then return 400
	e.GET("/api/v1/original").
		WithQuery("url", "invalid-url").
		Expect().
		Status(http.StatusBadRequest)

	// When valid URL passed but not found, then return 404
	e.GET("/api/v1/original").
		WithQuery("url", "https://host.ru/link/123").
		Expect().
		Status(http.StatusNotFound)

	// When valid URL passed and found, then return 200
	e.GET("/api/v1/original").
		WithQuery("url", "https://service.com/link/t92YuUGbw1WY4V2LvozcwRHdoB").
		WithHost("service.com").
		Expect().
		Status(http.StatusOK).
		JSON().Object().HasValue("destination", "https://example.com")
}

func testRedirect(e *httpexpect.Expect) {
	// When valid URL passed but not found, then return 404
	e.GET("/link/123").
		Expect().
		Status(http.StatusNotFound)

	// When URL found, then return 301 and redirect to original URL
	e.GET("/link/t92YuUGbw1WY4V2LvozcwRHdoB").
		Expect().
		Status(http.StatusMovedPermanently).
		Header("Location").
		IsEqual("https://example.com")
}

func initHandler(store map[string]string) http.Handler {
	storage := shortener.NewInMemoryStorage(store)
	shortenUseCase := shortener.NewShortenUseCase(storage)
	generateShortenUseCase := shortener.NewGenerateShortenUseCase(storage)
	getOriginalUseCase := shortener.NewGetOriginalUseCase(storage)
	getForRedirect := shortener.NewGetOriginalForRedirectUseCase(storage)

	shortenHandler := NewShortenHandler(
		shortenUseCase,
		generateShortenUseCase,
		getOriginalUseCase,
		getForRedirect,
		zerolog.Logger{},
	)

	handler := NewEcho(shortenHandler)

	return handler
}
