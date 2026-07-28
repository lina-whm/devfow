package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	client *redis.Client
}

func NewQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

func (q *Queue) Enqueue(ctx context.Context, queue string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal queue payload: %w", err)
	}
	return q.client.LPush(ctx, fmt.Sprintf("queue:%s", queue), data).Err()
}

func (q *Queue) Dequeue(ctx context.Context, queue string, timeout time.Duration) (string, error) {
	result, err := q.client.BRPop(ctx, timeout, fmt.Sprintf("queue:%s", queue)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("dequeue: %w", err)
	}
	if len(result) < 2 {
		return "", nil
	}
	return result[1], nil
}