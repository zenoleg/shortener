package app

import (
	"context"

	"github.com/zenoleg/shortener/internal/transport"
	"github.com/zenoleg/shortener/third_party/logger"
)

type App struct {
	server *transport.Server
}

func Init(appVersion string) *App {
	log := logger.NewLogger(logger.NewConfig(), appVersion)
	transportConfig := transport.NewConfig()

	server := transport.NewServer(transportConfig, log)

	return &App{server: server}
}

func (a *App) Start() error {
	return a.server.Run()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
