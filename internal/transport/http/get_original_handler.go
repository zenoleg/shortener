package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/usecase"
)

type (
	GetOriginalURLHandler struct {
		getOriginal usecase.GetOriginalUseCase
		logger      zerolog.Logger
	}

	GetOriginalRequest struct {
		URL string `query:"url"`
	}

	GetOriginalResponse struct {
		Destination string `json:"destination"`
	}
)

func NewGetOriginalURLHandler(getOriginal usecase.GetOriginalUseCase, logger zerolog.Logger) GetOriginalURLHandler {
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

	destination, err := h.getOriginal.Do(url)

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
