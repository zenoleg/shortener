package app

import (
	"context"

	"github.com/zenoleg/shortener/internal/shortener"
	"github.com/zenoleg/shortener/internal/transport"
	"github.com/zenoleg/shortener/third_party/logger"
)

type App struct {
	server *transport.Server
}

func Init(appVersion string) *App {
	log := logger.NewLogger(logger.NewConfig(), appVersion)
	transportConfig := transport.NewConfig()

	storage := shortener.NewInMemoryStorage(map[string]string{})
	shortenUseCase := shortener.NewShortenUseCase(storage)
	generateShortenUseCase := shortener.NewGenerateShortenUseCase()
	getOriginalUseCase := shortener.NewGetOriginalUseCase(storage)

	shortenHandler := transport.NewShortenHandler(shortenUseCase, generateShortenUseCase, getOriginalUseCase, log)
	echo := transport.NewEcho(shortenHandler)
	server := transport.NewServer(transportConfig, echo, log)

	return &App{server: server}
}

func (a *App) Start() error {
	return a.server.Run()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
