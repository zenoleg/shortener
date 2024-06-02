package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zenoleg/shortener/internal/shortener"
)

func TestShortenHandler_Shorten(t *testing.T) {
	t.Parallel()

	t.Run("When URL did not pass validation, then return 400", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage()
		useCase := shortener.NewShortenUseCase(storage)
		handler := NewShortenHandler(useCase, zerolog.Logger{})

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareRequest("{}")

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.Shorten(ectx)) {
			response := ErrorResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "url: cannot be blank.", response.Message)
		}
	})

	t.Run("When invalid URL passed, then return 400", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage()
		useCase := shortener.NewShortenUseCase(storage)
		handler := NewShortenHandler(useCase, zerolog.Logger{})

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareRequest(`{"url": "invalid-url"}`)

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.Shorten(ectx)) {
			response := ErrorResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "url 'invalid-url' is invalid", response.Message)
		}
	})

	t.Run("When valid URL passed, then return 201", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage()
		useCase := shortener.NewShortenUseCase(storage)
		handler := NewShortenHandler(useCase, zerolog.Logger{})

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareRequest(`{"url": "https://example.com"}`)

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.Shorten(ectx)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})
}

func prepareRequest(body string) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return request
}
