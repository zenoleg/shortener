package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zenoleg/shortener/internal/app"
)

var version = "unknown"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() error {
	application, err := app.Init(version)
	if err != nil {
		return err
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := application.Stop(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}

		close(stopped)
	}()

	if err := application.Start(); err != nil {
		log.Fatalf("Application startup error: %s", err)
	}

	<-stopped

	log.Printf("Bye! 👋")

	return nil
}
