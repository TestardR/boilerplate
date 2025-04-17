package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boilerplate/config"
	"boilerplate/internal/application"
	httpv1 "boilerplate/internal/infrastructure/api/http_v1"
	"boilerplate/internal/infrastructure/persistence/postgres"
	"boilerplate/internal/infrastructure/system"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cfg, err := config.FromEnv()
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.Level(cfg.LogLevel),
	}))

	logger.InfoContext(ctx, "Starting service")

	postgresClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create postgres client", "error", err)
		os.Exit(1)
	}

	migration, err := migrate.New(cfg.Postgres.SourceURL(), cfg.Postgres.DBURL())
	if err != nil {
		logger.ErrorContext(ctx, "failed to create migrator", "error", err)
		os.Exit(1)
	}
	defer migration.Close() //nolint:staticcheck

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.ErrorContext(ctx, "failed to run migrations", "error", err)
		os.Exit(1)
	}

	userStore := postgres.NewUserStore(postgresClient.DB())

	userService := application.NewUserService(
		userStore,
		userStore,
		system.NewClock(),
	)

	handler := httpv1.NewHandler(userService)

	httpServer := httpv1.NewHttServer(cfg.HTTP, logger, handler)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "failed to start http server", "error", err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		logger.InfoContext(ctx, "Stopping HTTP Server...")

		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "HTTP Server shutdown error", "error", err.Error())
		}

		err = postgresClient.Close()
		if err != nil {
			logger.ErrorContext(ctx, "failed to close postgres client", "error", err)
		}
	}()

	<-stop
	logger.Info("Exiting Service")
}
