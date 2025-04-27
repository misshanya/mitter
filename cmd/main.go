package main

import (
	"context"
	"fmt"
	"github.com/misshanya/mitter/internal/app"
	"github.com/misshanya/mitter/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

//	@title		Mitter
//	@version	1.0

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	cfg := config.NewConfig()
	server := app.NewApp(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go server.Start(ctx)

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("shutting down")
	if err := server.Stop(ctx); err != nil {
		slog.Error("failed to stop server", slog.Any("err", err))
	}
}
