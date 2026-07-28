package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		key := "rate_limit:" + c.ClientIP()
		now := time.Now().Unix()

		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-60, 10))
		countCmd := pipe.ZCard(ctx, key)
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
		pipe.Expire(ctx, key, time.Minute)

		if _, err := pipe.Exec(ctx); err != nil {
			c.Next()
			return
		}

		if countCmd.Val() >= int64(requestsPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
