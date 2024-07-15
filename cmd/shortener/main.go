package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/storage"
	"github.com/zenoleg/shortener/internal/transport/http"
	"github.com/zenoleg/shortener/internal/transport/http/handler"
	"github.com/zenoleg/shortener/internal/usecase"
	"github.com/zenoleg/shortener/third_party/logger"
)

var version = "unknown"

func main() {
	appLogger := logger.NewLogger(logger.NewConfig(), version)
	transportConfig := http.NewConfig()

	idGenerator := domain.NewBase62IDGenerator()
	store := storage.NewInMemoryStorage(make(map[string]string))

	shortenUseCase := usecase.NewShortenUseCase(store, idGenerator)
	generateShortenUseCase := usecase.NewGetShortenUseCase(store, idGenerator)
	getOriginalUseCase := usecase.NewGetOriginalUseCase(store)
	getOriginalForRedirect := usecase.NewGetOriginalForRedirectUseCase(store)

	shortenHandler := handler.NewShortenHandler(shortenUseCase, appLogger)
	generateShortenHandler := handler.NewGetShortURLHandler(generateShortenUseCase, appLogger)
	getOriginalHandler := handler.NewGetOriginalURLHandler(getOriginalUseCase, appLogger)
	redirectHandler := handler.NewRedirectHandler(getOriginalForRedirect, appLogger)

	echo := http.NewEcho(shortenHandler, generateShortenHandler, getOriginalHandler, redirectHandler)
	server := http.NewServer(transportConfig, echo, appLogger)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			appLogger.Printf("HTTP Server Shutdown Error: %v", err)
		}
	}()

	if err := server.Run(); err != nil {
		appLogger.Fatal().Err(err).Msg("server startup error")
	}

	appLogger.Printf("Bye! 👋")
}
