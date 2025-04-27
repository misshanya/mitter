package app

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/misshanya/mitter/docs"
	"github.com/misshanya/mitter/internal/api/handler"
	"github.com/misshanya/mitter/internal/config"
	"github.com/misshanya/mitter/internal/db"
	"github.com/misshanya/mitter/internal/db/sqlc/storage"
	myMiddleware "github.com/misshanya/mitter/internal/middleware"
	"github.com/misshanya/mitter/internal/repository"
	"github.com/misshanya/mitter/internal/service/auth"
	"github.com/misshanya/mitter/internal/service/user"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/echo-swagger"
	"log/slog"
	"os"
)

type App struct {
	cfg *config.Config
	e   *echo.Echo
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Start(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Addr,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		slog.Error("failed to connect to redis")
		os.Exit(1)
	}

	// Init db connection
	conn, err := initDB(ctx, a.cfg.Postgres.URL)
	if err != nil {
		slog.Error("failed to connect to database")
		os.Exit(1)
	}

	if err := db.Migrate(sql.OpenDB(stdlib.GetConnector(*conn.Config().ConnConfig))); err != nil {
		slog.Error("failed to migrate database")
		os.Exit(1)
	}

	// Init SQL queries
	queries := storage.New(conn)

	a.e = echo.New()
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.Logger())

	// Swagger
	a.e.GET("/swagger/*", echoSwagger.WrapHandler)

	apiGroup := a.e.Group("/api")
	v1Group := apiGroup.Group("/v1")

	// Repos
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(rdb)

	// Services
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(userRepo, authRepo)

	// Middlewares
	authMiddleware := myMiddleware.NewAuthMiddleware(authRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, authMiddleware.RequireAuth)

	// Groups
	userGroup := v1Group.Group("/user")
	authGroup := v1Group.Group("/auth")

	// Apply middlewares
	userGroup.Use(authMiddleware.RequireAuth)

	// Connect handlers
	userHandler.Routes(userGroup)
	authHandler.Routes(authGroup)

	a.e.Logger.Fatal(a.e.Start(a.cfg.Server.Addr))
}

func (a *App) Stop(ctx context.Context) error {
	return a.e.Shutdown(ctx)
}

func initDB(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = 100 // Max 100 connections

	return pool, nil
}
