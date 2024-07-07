package handler

import (
	"context"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=GetOriginalUseCase
type GetOriginalUseCase interface {
	Do(ctx context.Context, url domain.URL) (domain.URL, error)
}

type (
	GetOriginalURLHandler struct {
		getOriginal GetOriginalUseCase
		logger      zerolog.Logger
	}

	GetOriginalRequest struct {
		URL string `query:"url"`
	}

	GetOriginalResponse struct {
		Destination string `json:"destination"`
	}
)

func NewGetOriginalURLHandler(getOriginal GetOriginalUseCase, logger zerolog.Logger) GetOriginalURLHandler {
	return GetOriginalURLHandler{
		getOriginal: getOriginal,
		logger:      logger,
	}
}

func (h *GetOriginalURLHandler) Handle(ectx echo.Context) error {
	req := GetOriginalRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.NoContent(http.StatusBadRequest)
	}

	url, err := domain.NewURL(req.URL)
	if err != nil {
		return ectx.NoContent(http.StatusBadRequest)
	}

	destination, err := h.getOriginal.Do(ectx.Request().Context(), url)

	if errors.Is(err, storage.ErrURLNotFound) {
		return ectx.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return ectx.String(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return ectx.JSON(
		http.StatusOK,
		GetOriginalResponse{Destination: destination.String()},
	)
}

func (r GetOriginalRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required),
	)
}
