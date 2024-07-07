package handler

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/zenoleg/shortener/internal/usecase"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=ShortenUseCase
type ShortenUseCase interface {
	Do(ctx context.Context, query usecase.ShortenQuery) (usecase.DestinationURL, error)
}

type (
	ShortenHandler struct {
		shorten ShortenUseCase
		logger  zerolog.Logger
	}

	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Destination string `json:"destination"`
	}
)

func NewShortenHandler(shorten ShortenUseCase, logger zerolog.Logger) ShortenHandler {
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
		ectx.Request().Context(),
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

func (r ShortenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required),
	)
}
