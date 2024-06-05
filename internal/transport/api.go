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
		shorten            shortener.ShortenUseCase
		generateShortenURL shortener.GetShortUseCase
		getOriginal        shortener.GetOriginalUseCase
		logger             zerolog.Logger
	}
)

func NewShortenHandler(
	shorten shortener.ShortenUseCase,
	generateShortenURL shortener.GetShortUseCase,
	getOriginal shortener.GetOriginalUseCase,
	logger zerolog.Logger,
) ShortenHandler {
	return ShortenHandler{
		shorten:            shorten,
		generateShortenURL: generateShortenURL,
		getOriginal:        getOriginal,
		logger:             logger,
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
		shortener.NewGenerateShortenQuery(
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

	destinationURL, err := h.getOriginal.Handle(req.ShortURL)

	if err != nil {
		h.logger.Err(err).Msg("failed to shorten url")

		return ectx.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err.Error()),
		)
	}

	return ectx.JSON(http.StatusOK, NewShortenResponse(destinationURL))
}
