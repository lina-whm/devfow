package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/devflow/devflow-backend/internal/api/handler"
	"github.com/devflow/devflow-backend/internal/api/router"
	"github.com/devflow/devflow-backend/internal/application/auth"
	"github.com/devflow/devflow-backend/internal/application/comment"
	"github.com/devflow/devflow-backend/internal/application/organization"
	"github.com/devflow/devflow-backend/internal/application/project"
	"github.com/devflow/devflow-backend/internal/application/task"
	"github.com/devflow/devflow-backend/internal/application/team"
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

	_, _ = maxprocs.Set()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.Database.URL, cfg.Database.MaxConns, cfg.Database.MinConns)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	rdb, err := redis.NewClient(ctx, cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	ginMode := cfg.Server.Mode
	if ginMode == "" || ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	userRepo := postgres.NewUserRepository(pool)

	authSvc := auth.NewService(userRepo)

	orgRepo := postgres.NewOrgRepository(pool)
	orgMemRepo := postgres.NewOrgMemberRepository(pool)
	invRepo := postgres.NewInvitationRepository(pool)
	orgSvc := organization.NewService(orgRepo, orgMemRepo, invRepo, userRepo)
	orgH := handler.NewOrgHandler(orgSvc)

	teamRepo := postgres.NewTeamRepository(pool)
	teamMemRepo := postgres.NewTeamMemberRepository(pool)
	_ = team.NewService(teamRepo, teamMemRepo)

	projectRepo := postgres.NewProjectRepository(pool)
	projectSvc := project.NewService(projectRepo)
	projectH := handler.NewProjectHandler(projectSvc)

	_ = postgres.NewBoardRepository(pool)
	_ = postgres.NewColumnRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	taskSvc := task.NewService(taskRepo, userRepo)
	taskH := handler.NewTaskHandler(taskSvc)

	commentRepo := postgres.NewCommentRepository(pool)
	commentSvc := comment.NewService(commentRepo)
	commentH := handler.NewCommentHandler(commentSvc)

	healthH := handler.NewHealthHandler(pool, rdb)
	authH := handler.NewAuthHandler(authSvc, rdb, cfg.JWT.Secret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	r := router.Setup(authH, healthH, taskH, orgH, projectH, commentH, cfg.JWT.Secret)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	middleware.GracefulShutdown(ctx, cancel, 30*time.Second, func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	})
}
