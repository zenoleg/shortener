package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/shortener"
)

type (
	ShortenHandler struct {
		shorten            shortener.ShortenUseCase
		generateShortenURL shortener.GenerateShortenUseCase
		logger             zerolog.Logger
	}
)

func NewShortenHandler(shorten shortener.ShortenUseCase, logger zerolog.Logger) ShortenHandler {
	return ShortenHandler{
		shorten: shorten,
		logger:  logger,
	}
}

func (h *ShortenHandler) Shorten(ectx echo.Context) error {
	req := ShortenRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	err = h.shorten.Handle(req.URL)
	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.NoContent(http.StatusCreated)
}

func (h *ShortenHandler) GenerateShortenURL(ectx echo.Context) error {
	req := ShortenRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	destinationURL, err := h.generateShortenURL.Handle(
		shortener.NewGenerateShortenQuery(
			ectx.Scheme() == "https",
			ectx.Request().Host,
			req.URL,
		),
	)

	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.JSON(http.StatusOK, NewShortenResponse(destinationURL.String()))
}
