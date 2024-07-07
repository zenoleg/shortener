package app

import (
	"context"

	"github.com/zenoleg/shortener/internal/infrastracture"
	"github.com/zenoleg/shortener/internal/shortener"
	"github.com/zenoleg/shortener/internal/transport"
	"github.com/zenoleg/shortener/internal/transport/http"
	"github.com/zenoleg/shortener/third_party/logger"
)

type App struct {
	server  *transport.Server
	cleanup func()
}

func Init(appVersion string) (*App, error) {
	log := logger.NewLogger(logger.NewConfig(), appVersion)
	transportConfig := transport.NewConfig()

	levelDBConfig := infrastracture.NewConfig()
	levelDBConnection, closeDB, err := infrastracture.NewLevelDBConnection(levelDBConfig, log)
	if err != nil {
		return nil, err
	}

	storage := shortener.NewLevelDBStorage(levelDBConnection, log)
	shortenUseCase := shortener.NewShortenUseCase(storage)
	generateShortenUseCase := shortener.NewGenerateShortenUseCase(storage)
	getOriginalUseCase := shortener.NewGetOriginalUseCase(storage)
	getOriginalForRedirect := shortener.NewGetOriginalForRedirectUseCase(storage)

	shortenHandler := transport.NewShortenHandler(
		shortenUseCase,
		generateShortenUseCase,
		getOriginalUseCase,
		getOriginalForRedirect,
		log,
	)
	echo := http.NewEcho(shortenHandler)
	server := transport.NewServer(transportConfig, echo, log)

	return &App{server: server, cleanup: closeDB}, nil
}

func (a *App) Start() error {
	return a.server.Run()
}

func (a *App) Stop(ctx context.Context) error {
	defer a.cleanup()

	return a.server.Shutdown(ctx)
}
