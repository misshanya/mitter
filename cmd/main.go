package main

import (
	"context"
	"github.com/misshanya/mitter/internal/app"
	"github.com/misshanya/mitter/internal/config"
)

func main() {
	cfg := config.NewConfig()
	server := app.NewApp(cfg)

	ctx := context.Background()
	server.Start(ctx)
}
