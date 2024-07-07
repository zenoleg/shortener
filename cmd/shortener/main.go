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
	application, err := app.Init(version)
	if err != nil {
		panic(err)
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := application.Stop(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
	}()

	if err = application.Start(); err != nil {
		log.Fatalf("Application startup error: %s", err)
	}

	log.Printf("Bye! 👋")
}
