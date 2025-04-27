package main

import (
	"context"
	"github.com/misshanya/mitter/internal/app"
	"github.com/misshanya/mitter/internal/config"
)

//	@title		Mitter
//	@version	1.0

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	cfg := config.NewConfig()
	server := app.NewApp(cfg)

	ctx := context.Background()
	server.Start(ctx)
}
