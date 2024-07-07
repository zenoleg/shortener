package handler

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=GetOriginalForRedirect
type GetOriginalForRedirect interface {
	Do(domain.ID) (domain.URL, error)
}

type (
	RedirectHandler struct {
		getForRedirect GetOriginalForRedirect
		logger         zerolog.Logger
	}

	RedirectRequest struct {
		ID string `param:"id"`
	}
)

func NewRedirectHandler(getForRedirect GetOriginalForRedirect, logger zerolog.Logger) RedirectHandler {
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
	}

	if err != nil {
		return ectx.NoContent(http.StatusInternalServerError)
	}

	return ectx.Redirect(http.StatusMovedPermanently, destinationURL.String())
}

func (r RedirectRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}
