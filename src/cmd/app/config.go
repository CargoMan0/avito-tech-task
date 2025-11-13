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
}

type HTTPServer struct {
	ListenAddr      string        `env:"LISTEN_ADDR,notEmpty"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,notEmpty"`
}

type Database struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty"`
	User     string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Name     string `env:"NAME,notEmpty"`
	SSLMode  string `env:"SSL_MODE,notEmpty"`
}

type Logger struct {
	LogLevel int  `env:"LOG_LEVEL,notEmpty"`
	IsJSON   bool `env:"IS_JSON,notEmpty"`
}

type Migrations struct {
	Dir     string `env:"DIR,notEmpty"`
	Enabled bool   `env:"ENABLED" envDefault:"false"`
}

func loadConfigFromEnv() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return c, nil
}
