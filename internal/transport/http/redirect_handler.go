package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
	"github.com/zenoleg/shortener/internal/usecase"
)

type (
	RedirectHandler struct {
		getForRedirect usecase.GetOriginalForRedirectUseCase
		logger         zerolog.Logger
	}

	RedirectRequest struct {
		ID string `param:"shortID"`
	}
)

func NewRedirectHandler(getForRedirect usecase.GetOriginalForRedirectUseCase, logger zerolog.Logger) RedirectHandler {
	return RedirectHandler{
		getForRedirect: getForRedirect,
		logger:         logger,
	}
}

func (h *RedirectHandler) Handle(ectx echo.Context) error {
	req := RedirectRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.NoContent(http.StatusBadRequest)
	}

	destinationURL, err := h.getForRedirect.Do(domain.ID(req.ID))

	if errors.Is(err, storage.ErrURLNotFound) {
		return ectx.NoContent(http.StatusNotFound)
	} else if err != nil {
		return ectx.NoContent(http.StatusInternalServerError)
	}

	return ectx.Redirect(http.StatusMovedPermanently, destinationURL.String())
}
