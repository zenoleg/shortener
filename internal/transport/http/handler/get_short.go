package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/usecase"
)

type (
	GetShortURLHandler struct {
		getShort usecase.GetShortUseCase
		logger   zerolog.Logger
	}

	GetShortURLRequest struct {
		URL string `query:"url"`
	}

	GetShortURLResponse struct {
		Destination string `json:"destination"`
	}
)

func NewGetShortURLHandler(getShort usecase.GetShortUseCase, logger zerolog.Logger) GetShortURLHandler {
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
