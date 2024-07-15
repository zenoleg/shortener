package handler

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/zenoleg/shortener/internal/usecase"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=GetShort
type GetShort interface {
	Do(ctx context.Context, query usecase.GetShortURLQuery) (usecase.DestinationURL, error)
}

type (
	GetShortURLHandler struct {
		getShort GetShort
		logger   zerolog.Logger
	}

	GetShortURLRequest struct {
		URL string `query:"url"`
	}

	GetShortURLResponse struct {
		Destination string `json:"destination"`
	}
)

func NewGetShortURLHandler(getShort GetShort, logger zerolog.Logger) GetShortURLHandler {
	return GetShortURLHandler{
		getShort: getShort,
		logger:   logger,
	}
}

func (h *GetShortURLHandler) Handle(ectx echo.Context) error {
	req := GetShortURLRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.NoContent(http.StatusBadRequest)
	}

	destination, err := h.getShort.Do(
		ectx.Request().Context(),
		usecase.NewGetShortURLQuery(
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
		GetShortURLResponse{Destination: destination.String()},
	)
}

func (r GetShortURLRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required),
	)
}
