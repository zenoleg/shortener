package transport

import (
	"context"
	"net/http"

	"emperror.dev/errors"
	"github.com/rs/zerolog"
)

type Server struct {
	srv    *http.Server
	logger zerolog.Logger
}

func NewServer(cfg Config, logger zerolog.Logger) *Server {
	return &Server{
		srv: &http.Server{
			Addr: cfg.Address,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info().Msg("starting server")
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Info().Err(err).Msg("server stopped")

		return err
	}

	s.logger.Info().Msg("server stopped")

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("server returned an error")
	}

	return nil
}
