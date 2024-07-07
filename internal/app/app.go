package app

import (
	"context"

	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
	"github.com/zenoleg/shortener/internal/transport/http"
	"github.com/zenoleg/shortener/internal/transport/http/handler"
	"github.com/zenoleg/shortener/internal/usecase"
	"github.com/zenoleg/shortener/third_party/logger"
)

type App struct {
	server *http.Server
}

func Init(appVersion string) (*App, error) {
	log := logger.NewLogger(logger.NewConfig(), appVersion)
	transportConfig := http.NewConfig()

	idGenerator := domain.NewBase62IDGenerator()
	store := storage.NewInMemoryStorage(make(map[string]string))

	shortenUseCase := usecase.NewShortenUseCase(store, idGenerator)
	generateShortenUseCase := usecase.NewGetShortenUseCase(store, idGenerator)
	getOriginalUseCase := usecase.NewGetOriginalUseCase(store)
	getOriginalForRedirect := usecase.NewGetOriginalForRedirectUseCase(store)

	shortenHandler := handler.NewShortenHandler(shortenUseCase, log)
	generateShortenHandler := handler.NewGetShortURLHandler(generateShortenUseCase, log)
	getOriginalHandler := handler.NewGetOriginalURLHandler(getOriginalUseCase, log)
	redirectHandler := handler.NewRedirectHandler(getOriginalForRedirect, log)

	echo := http.NewEcho(shortenHandler, generateShortenHandler, getOriginalHandler, redirectHandler)
	server := http.NewServer(transportConfig, echo, log)

	return &App{server: server}, nil
}

func (a *App) Start() error {
	return a.server.Run()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
