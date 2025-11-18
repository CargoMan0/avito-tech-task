package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	HTTPServer HTTPServer `envPrefix:"HTTP_SERVER_"`
	Database   Database   `envPrefix:"DATABASE_"`
	Logger     Logger     `envPrefix:"LOGGER_"`
	Migrations Migrations `envPrefix:"MIGRATIONS_"`
	Auth       Auth       `envPrefix:"AUTH_"`
}

type HTTPServer struct {
	ListenAddr      string        `env:"LISTEN_ADDR" envDefault:"0.0.0.0:3000"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type Database struct {
	Host     string `env:"HOST" envDefault:"db"`
	Port     int    `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"dev"`
	Password string `env:"PASSWORD" envDefault:"dev"`
	Name     string `env:"NAME" envDefault:"dev"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

type Logger struct {
	LogLevel int  `env:"LOG_LEVEL" envDefault:"-4"`
	IsJSON   bool `env:"IS_JSON" envDefault:"false"`
}

type Migrations struct {
	Dir     string `env:"DIR" envDefault:"./migrations"`
	Enabled bool   `env:"ENABLED" envDefault:"true"`
}

type Auth struct {
	Secret string `env:"SECRET" envDefault:"ultra-secret-jwt"`
}

func loadConfigFromEnv() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return c, nil
}
