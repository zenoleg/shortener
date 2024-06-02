package transport

import (
	"context"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Server struct {
	echo   *echo.Echo
	logger zerolog.Logger
	cfg    Config
}

func NewServer(cfg Config, echoServer *echo.Echo, logger zerolog.Logger) *Server {
	return &Server{
		echo:   echoServer,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Server) Run() error {
	s.logger.Info().Msg("starting server")
	err := s.echo.Start(s.cfg.Address)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Info().Err(err).Msg("server stopped")

		return err
	}

	s.logger.Info().Msg("server stopped")

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("server returned an error")
	}

	return nil
}
