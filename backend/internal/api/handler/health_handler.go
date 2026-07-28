package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func NewHealthHandler(pool *pgxpool.Pool, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{pool: pool, rdb: rdb}
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx := context.Background()
	checks := gin.H{"status": "ok"}

	if err := h.pool.Ping(ctx); err != nil {
		checks["status"] = "degraded"
		checks["database"] = "unreachable"
	} else {
		checks["database"] = "ok"
	}

	if _, err := h.rdb.Ping(ctx).Result(); err != nil {
		checks["status"] = "degraded"
		checks["redis"] = "unreachable"
	} else {
		checks["redis"] = "ok"
	}

	status := http.StatusOK
	if checks["status"] != "ok" {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, checks)
}
