package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zenoleg/shortener/internal/shortener"
)

func TestShortenHandler_Shorten(t *testing.T) {
	t.Parallel()

	t.Run("When URL did not pass validation, then return 400", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage(map[string]string{})
		useCase := shortener.NewShortenUseCase(storage)
		handler := ShortenHandler{shorten: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := preparePostRequest("{}")

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
		storage := shortener.NewInMemoryStorage(map[string]string{})
		useCase := shortener.NewShortenUseCase(storage)
		handler := ShortenHandler{shorten: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := preparePostRequest(`{"url": "invalid-url"}`)

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
		storage := shortener.NewInMemoryStorage(map[string]string{})
		useCase := shortener.NewShortenUseCase(storage)
		handler := ShortenHandler{shorten: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := preparePostRequest(`{"url": "https://example.com"}`)

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.Shorten(ectx)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})
}

func TestShortenHandler_GenerateShortenURL(t *testing.T) {
	t.Parallel()

	t.Run("When URL did not pass validation, then return 400", func(t *testing.T) {
		useCase := shortener.NewGenerateShortenUseCase()
		handler := ShortenHandler{generateShortenURL: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareGetRequest("")

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

	t.Run("When valid URL passed, then return 200 and short URL", func(t *testing.T) {
		useCase := shortener.NewGenerateShortenUseCase()
		handler := ShortenHandler{generateShortenURL: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareGetRequest("?url=https://example.com")

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.GenerateShortenURL(ectx)) {
			resp := DestinationResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "http://example.com/link/t92YuUGbw1WY4V2LvozcwRHdoB", resp.Destination)
		}
	})
}

func TestShortenHandler_GetOriginal(t *testing.T) {
	t.Parallel()

	t.Run("When URL did not pass validation, then return 400", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage(map[string]string{})
		useCase := shortener.NewGetOriginalUseCase(storage)
		handler := ShortenHandler{getOriginal: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareGetRequest("")

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/original")

		if assert.NoError(t, handler.GetOriginal(ectx)) {
			response := ErrorResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ShortURL: cannot be blank.", response.Message)
		}
	})

	t.Run("When invalid URL passed, then return 200 and original URL", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage(map[string]string{
			"t92YuUGbw1WY4V2LvozcwRHdoB": "http://example.com",
		})

		useCase := shortener.NewGetOriginalUseCase(storage)
		handler := ShortenHandler{getOriginal: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareGetRequest("?short_url=http://example.com/t92YuUGbw1WY4V2LvozcwRHdoB")

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.GetOriginal(ectx)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("When valid URL passed, then return 200 and original URL", func(t *testing.T) {
		storage := shortener.NewInMemoryStorage(map[string]string{
			"t92YuUGbw1WY4V2LvozcwRHdoB": "http://example.com",
		})

		useCase := shortener.NewGetOriginalUseCase(storage)
		handler := ShortenHandler{getOriginal: useCase}

		e := echo.New()
		e.Binder = NewValidationBinder()

		request := prepareGetRequest("?short_url=http://example.com/link/t92YuUGbw1WY4V2LvozcwRHdoB")

		rec := httptest.NewRecorder()
		ectx := e.NewContext(request, rec)
		ectx.SetPath("/api/v1/shorten")

		if assert.NoError(t, handler.GetOriginal(ectx)) {
			resp := DestinationResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "http://example.com", resp.Destination)
		}
	})
}

func preparePostRequest(body string) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return request
}

func prepareGetRequest(queryParams string) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/"+queryParams, nil)

	return request
}
