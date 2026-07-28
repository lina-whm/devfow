package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		requestID, _ := c.Get("request_id")

		slog.Info("request",
			"method", method,
			"path", path,
			"status", status,
			"duration", duration,
			"request_id", requestID,
		)
	}
}
