package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/http/handlers/routes"
	"github.com/CargoMan0/avito-tech-task/internal/repository/impl"
	"github.com/CargoMan0/avito-tech-task/internal/service"
	"github.com/CargoMan0/avito-tech-task/pkg/database"
	"github.com/CargoMan0/avito-tech-task/pkg/migrations"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run() returned error: %v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)
	defer cancel()

	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := getLogger(cfg.Logger.LogLevel, cfg.Logger.IsJSON)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	logger.Info("Logger setup successfully")

	logger.Info("Connecting to SQL database")
	sqlDB, err := newSQLDatabase(cfg.Database)
	if err != nil {
		return fmt.Errorf("create new sql database: %w", err)
	}
	defer func() {
		logger.Info("Closing SQL database")
		closeErr := sqlDB.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close sql database: %w", closeErr))
			return
		}

		logger.Info("SQL database closed")
	}()

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	logger.Info("Connected to SQL database")

	if cfg.Migrations.Enabled {
		logger.Info("Migrations enabled - running migrations")
		err = migrations.Run(ctx, sqlDB, cfg.Migrations.Dir)
		if err != nil {
			return fmt.Errorf("run migrations: %w", err)
		}
		logger.Info("Successfully run migrations")
	} else {
		logger.Info("Migrations disabled - skipping migrations")
	}

	// Repositories
	pullRequestsRepo := impl.NewPullRequestRepository(sqlDB)
	usersRepo := impl.NewUserRepository(sqlDB)
	teamsRepo := impl.NewTeamRepository(sqlDB)

	// Service
	srvc := service.NewService(
		pullRequestsRepo,
		usersRepo,
		teamsRepo,
	)

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, srvc)

	srv := &http.Server{
		Addr:    cfg.HTTPServer.ListenAddr,
		Handler: mux,
	}

	errChan := make(chan error, 2)
	stopWg := &sync.WaitGroup{}

	stopWg.Add(1)
	go func() {
		defer stopWg.Done()

		logger.Info("Starting HTTP server")
		listenErr := srv.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			errChan <- fmt.Errorf("listen and serve: %v", listenErr)
			cancel()
		}
	}()

	stopWg.Add(1)
	go func(cfg HTTPServer) {
		defer stopWg.Done()

		<-ctx.Done()
		logger.Info("Shutdown signal received. Shutting down HTTP server")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		shutdownErr := srv.Shutdown(timeoutCtx)
		if shutdownErr != nil && !errors.Is(shutdownErr, http.ErrServerClosed) {
			errChan <- fmt.Errorf("shutdown http server: %v", shutdownErr)
		}

		logger.Info("HTTP server shutdown successfully")
	}(cfg.HTTPServer)

	go func() {
		stopWg.Wait()
		close(errChan)
	}()

	for errFromGoroutine := range errChan {
		err = errors.Join(err, errFromGoroutine)
	}

	return
}

func newSQLDatabase(cfg Database) (*sql.DB, error) {
	dbCfg := database.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.Name,
		SSLMode:  cfg.SSLMode,
	}

	return database.NewSQL(dbCfg)
}

func getLogger(level int, isJson bool) (*slog.Logger, error) {
	var (
		handler     slog.Handler
		handlerOpts = &slog.HandlerOptions{
			Level: slog.Level(level),
		}
	)

	if isJson {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	logger := slog.New(handler)
	return logger, nil
}
