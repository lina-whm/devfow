package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/devflow/devflow-backend/internal/infrastructure/postgres"
	"github.com/devflow/devflow-backend/internal/infrastructure/redis"
	"github.com/devflow/devflow-backend/internal/pkg/config"
	"github.com/devflow/devflow-backend/internal/pkg/logger"
	"github.com/devflow/devflow-backend/internal/pkg/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Logger.Level, cfg.Logger.Format, nil)
	slog.SetDefault(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.NewPool(ctx, cfg.Database.URL, cfg.Database.MaxConns, cfg.Database.MinConns)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	rdb, err := redis.NewClient(ctx, cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer func() { _ = rdb.Close() }()

	slog.Info("worker started")

	middleware.GracefulShutdown(ctx, cancel, 30*time.Second, func(ctx context.Context) error {
		slog.Info("worker shutting down")
		return nil
	})
}
