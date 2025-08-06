package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   server   `env:"SERVER" env-required:"true"`
	Redis    redis    `env:"REDIS" env-required:"true"`
	Postgres postgres `env:"POSTGRES" env-required:"true"`
	Mode     string   `env:"MODE" env-default:"PROD"`
}

type server struct {
	Addr string `env:"SERVER_ADDRESS" env-required:"true"`
}

type redis struct {
	Addr     string `env:"REDIS_ADDRESS" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB" env-required:"true"`
}

type postgres struct {
	URL      string `env:"PG_URL" env-required:"true"`
	MaxConns int32  `env:"PG_MAX_CONNS" env-default:"100"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
