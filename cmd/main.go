package main

import (
	"context"
	"github.com/misshanya/mitter/internal/app"
	"github.com/misshanya/mitter/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title		Mitter
//	@version	1.0

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	logger := setupLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("failed to read config", "error")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to create app", "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error)
	go server.Start(errChan)

	select {
	case err := <-errChan:
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop server", "err", err)
			os.Exit(1)
		}
	}
}

func setupLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	logger := slog.New(handler)
	return logger
}
