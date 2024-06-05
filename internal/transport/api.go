package transport

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zenoleg/shortener/internal/shortener"
)

type (
	ShortenHandler struct {
		shorten                shortener.ShortenUseCase
		generateShortenURL     shortener.GetShortUseCase
		getOriginal            shortener.GetOriginalUseCase
		getOriginalForRedirect shortener.GetOriginalForRedirectUseCase
		logger                 zerolog.Logger
	}
)

func NewShortenHandler(
	shorten shortener.ShortenUseCase,
	generateShortenURL shortener.GetShortUseCase,
	getOriginal shortener.GetOriginalUseCase,
	getOriginalForRedirect shortener.GetOriginalForRedirectUseCase,
	logger zerolog.Logger,
) ShortenHandler {
	return ShortenHandler{
		shorten:                shorten,
		generateShortenURL:     generateShortenURL,
		getOriginal:            getOriginal,
		getOriginalForRedirect: getOriginalForRedirect,
		logger:                 logger,
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

	destination, err := h.shorten.Handle(
		shortener.NewShortenQuery(
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

	return ectx.JSON(http.StatusCreated, NewShortenResponse(destination.String()))
}

func (h *ShortenHandler) GetShortURL(ectx echo.Context) error {
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
		shortener.NewShortenQuery(
			ectx.Scheme() == "https",
			ectx.Request().Host,
			req.URL,
		),
	)

	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		if errors.Is(err, shortener.ErrNotFound) {
			return ectx.JSON(
				http.StatusNotFound,
				NewErrorResponse(err.Error()),
			)
		}

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.JSON(http.StatusOK, NewShortenResponse(destinationURL.String()))
}

func (h *ShortenHandler) GetOriginal(ectx echo.Context) error {
	req := OriginalRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	destinationURL, err := h.getOriginal.Handle(req.URL)

	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		if errors.Is(err, shortener.ErrNotFound) {
			return ectx.JSON(
				http.StatusNotFound,
				NewErrorResponse(err.Error()),
			)
		}

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.JSON(http.StatusOK, NewShortenResponse(destinationURL))
}

func (h *ShortenHandler) Redirect(ectx echo.Context) error {
	req := RedirectRequest{}

	err := ectx.Bind(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to bind request")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	destinationURL, err := h.getOriginalForRedirect.Handle(req.ShortID)

	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		if errors.Is(err, shortener.ErrNotFound) {
			return ectx.JSON(
				http.StatusNotFound,
				NewErrorResponse(err.Error()),
			)
		}

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.Redirect(http.StatusMovedPermanently, destinationURL)
}
