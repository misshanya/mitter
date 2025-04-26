package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Server   server
	Redis    redis
	Postgres postgres
}

type server struct {
	Addr string
}

type redis struct {
	Addr     string
	Password string
	DB       int
}

type postgres struct {
	URL string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file")
	}

	// server
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:8080"
	}

	// redis
	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	var (
		redisDB int
		err     error
	)
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr != "" {
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			slog.Error("Error parsing REDIS_DB")
			os.Exit(1)
		}
	}

	// postgres
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		slog.Error("Error parsing PG_URL")
		os.Exit(1)
	}

	return &Config{
		Server: server{
			Addr: serverAddr,
		},
		Redis: redis{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		Postgres: postgres{
			URL: pgURL,
		},
	}
}
