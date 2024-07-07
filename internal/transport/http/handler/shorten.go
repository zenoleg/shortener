package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/usecase"
)

type (
	ShortenHandler struct {
		shorten usecase.ShortenUseCase
		logger  zerolog.Logger
	}

	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Destination string `json:"destination"`
	}
)

func NewShortenHandler(shorten usecase.ShortenUseCase, logger zerolog.Logger) ShortenHandler {
	return ShortenHandler{
		shorten: shorten,
		logger:  logger,
	}
}

func (h *ShortenHandler) Handle(ectx echo.Context) error {
	req := ShortenRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.NoContent(http.StatusBadRequest)
	}

	destination, err := h.shorten.Do(
		usecase.NewShortenQuery(
			ectx.Scheme() == "https",
			ectx.Request().Host,
			req.URL,
		),
	)

	if err != nil {
		return ectx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return ectx.JSON(
		http.StatusOK,
		ShortenResponse{Destination: destination.String()},
	)
}
