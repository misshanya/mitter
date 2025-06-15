package app

import (
	"context"
	"database/sql"
	"errors"
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

func NewApp(cfg *config.Config, logger *slog.Logger) *App {
	return &App{cfg: cfg, l: logger}
}

func (a *App) Start(ctx context.Context, errChan chan<- error) {
	a.rdb = redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Addr,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	})

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		errChan <- err
		return
	}

	// Init db connection
	a.dbPool, err = initDB(ctx, a.cfg.Postgres.URL)
	if err != nil {
		errChan <- err
		return
	}

	if err := db.Migrate(sql.OpenDB(stdlib.GetConnector(*a.dbPool.Config().ConnConfig))); err != nil {
		errChan <- err
		return
	}

	// Init SQL queries
	queries := storage.New(a.dbPool)

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

	// Swagger
	a.e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Prometheus
	a.e.Use(echoprometheus.NewMiddleware("mitter"))
	a.e.GET("/metrics", echoprometheus.NewHandler())

	// Custom metrics
	userMetrics := metrics.NewUserMetrics()
	mittMetrics := metrics.NewMittMetrics()

	apiGroup := a.e.Group("/api")
	v1Group := apiGroup.Group("/v1")

	// Repos
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(a.rdb)
	mittRepo := repository.NewMittRepository(queries)

	// Services
	userService := user.NewUserService(userRepo, userMetrics, a.l)
	authService := auth.NewAuthService(userRepo, authRepo, userMetrics, a.l)
	mittService := mitt.NewService(mittRepo, mittMetrics, userRepo, a.l)

	// Middlewares
	authMiddleware := myMiddleware.NewAuthMiddleware(authRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, authMiddleware.RequireAuth)
	mittHandler := handler.NewMittHandler(mittService, authMiddleware.RequireAuth)

	// Groups
	userGroup := v1Group.Group("/user")
	authGroup := v1Group.Group("/auth")
	mittGroup := v1Group.Group("/mitt")

	// Apply middlewares
	userGroup.Use(authMiddleware.RequireAuth)

	// Connect handlers
	userHandler.Routes(userGroup)
	authHandler.Routes(authGroup)
	mittHandler.Routes(mittGroup)

	if err := a.e.Start(a.cfg.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.l.Info("[!] Shutting down...")

	var stopErr error

	// Stop server
	a.l.Info("Stopping http server...")
	if err := a.e.Shutdown(ctx); err != nil {
		a.l.Error("failed to stop http server", slog.Any("error", err))
		stopErr = errors.Join(stopErr, err)
	}

	// Close DB pool
	a.l.Info("Closing database pool...")
	a.dbPool.Close()

	// Close Redis connection
	a.l.Info("Closing Redis connection...")
	if err := a.rdb.Close(); err != nil {
		a.l.Error("failed to close redis connection", slog.Any("error", err))
		stopErr = errors.Join(stopErr, err)
	}

	if stopErr != nil {
		return stopErr
	}

	a.l.Info("Stopped gracefully")
	return nil
}

func initDB(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = 100 // Max 100 connections

	return pool, nil
}
