package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/misshanya/mitter/docs"
	"github.com/misshanya/mitter/internal/api/handler"
	"github.com/misshanya/mitter/internal/config"
	"github.com/misshanya/mitter/internal/db"
	"github.com/misshanya/mitter/internal/db/sqlc/storage"
	"github.com/misshanya/mitter/internal/metrics"
	myMiddleware "github.com/misshanya/mitter/internal/middleware"
	"github.com/misshanya/mitter/internal/repository"
	"github.com/misshanya/mitter/internal/service/auth"
	"github.com/misshanya/mitter/internal/service/mitt"
	"github.com/misshanya/mitter/internal/service/user"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log/slog"
	"net/http"
)

type App struct {
	cfg    *config.Config
	e      *echo.Echo
	l      *slog.Logger
	dbPool *pgxpool.Pool
	rdb    *redis.Client
}

// New creates and initializes a new instance of App
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	a := &App{cfg: cfg, l: logger}

	if err := a.initRedis(ctx); err != nil {
		return nil, err
	}

	if err := a.initDB(ctx); err != nil {
		return nil, err
	}

	queries := storage.New(a.dbPool)

	a.initEcho()

	userMetrics := metrics.NewUserMetrics()
	mittMetrics := metrics.NewMittMetrics()

	apiGroup := a.e.Group("/api")
	v1Group := apiGroup.Group("/v1")

	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(a.rdb)
	mittRepo := repository.NewMittRepository(queries)

	userService := user.NewUserService(userRepo, userMetrics, a.l)
	authService := auth.NewAuthService(userRepo, authRepo, userMetrics, a.l)
	mittService := mitt.NewService(mittRepo, mittMetrics, userRepo, a.l)

	authMiddleware := myMiddleware.NewAuthMiddleware(authRepo)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, authMiddleware.RequireAuth)
	mittHandler := handler.NewMittHandler(mittService, authMiddleware.RequireAuth)

	userGroup := v1Group.Group("/user")
	authGroup := v1Group.Group("/auth")
	mittGroup := v1Group.Group("/mitt")

	userGroup.Use(authMiddleware.RequireAuth)

	userHandler.Routes(userGroup)
	authHandler.Routes(authGroup)
	mittHandler.Routes(mittGroup)

	return a, nil
}

// Start performs start for http server
func (a *App) Start(errChan chan<- error) {
	if err := a.e.Start(a.cfg.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
}

// Stop performs a graceful shutdown for all components
func (a *App) Stop(ctx context.Context) error {
	a.l.Info("[!] Shutting down...")

	var stopErr error

	a.l.Info("Stopping http server...")
	if err := a.e.Shutdown(ctx); err != nil {
		stopErr = errors.Join(stopErr, err)
	}

	a.l.Info("Closing database pool...")
	a.dbPool.Close()

	a.l.Info("Closing Redis connection...")
	if err := a.rdb.Close(); err != nil {
		stopErr = errors.Join(stopErr, err)
	}

	if stopErr != nil {
		return stopErr
	}

	a.l.Info("Stopped gracefully")
	return nil
}

// initRedis creates a Redis client and tries to ping server
func (a *App) initRedis(ctx context.Context) error {
	a.rdb = redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Addr,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	})

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}

// initDB initializes a new PostgreSQL pool and migrates schemas
func (a *App) initDB(ctx context.Context) error {
	dbPool, err := initDB(ctx, a.cfg.Postgres.URL, a.cfg.Postgres.MaxConns)
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}
	a.dbPool = dbPool

	if err := db.Migrate(sql.OpenDB(stdlib.GetConnector(*a.dbPool.Config().ConnConfig))); err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	return nil
}

// initEcho sets up a new Echo with recoverer, logger, CORS middlewares & Swagger and Prometheus handlers
func (a *App) initEcho() {
	a.e = echo.New()
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.Logger())

	if a.cfg.Mode == "DEV" {
		a.l.Info("[!] Running in DEV mode")
		a.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	a.e.GET("/swagger/*", echoSwagger.WrapHandler)

	a.e.Use(echoprometheus.NewMiddleware("mitter"))
	a.e.GET("/metrics", echoprometheus.NewHandler())
}

func initDB(ctx context.Context, dbURL string, maxConns int32) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = maxConns

	return pool, nil
}
