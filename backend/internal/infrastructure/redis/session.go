package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{client: client}
}

func (s *SessionStore) StoreRefreshToken(ctx context.Context, tokenHash, userID, orgID, role string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s", tokenHash)
	data := map[string]string{"user_id": userID, "org_id": orgID, "role": role}
	return s.client.HSet(ctx, key, data).Err()
}

func (s *SessionStore) GetRefreshTokenData(ctx context.Context, tokenHash string) (userID, orgID, role string, err error) {
	key := fmt.Sprintf("refresh:%s", tokenHash)
	data, err := s.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return "", "", "", fmt.Errorf("refresh token not found")
	}
	return data["user_id"], data["org_id"], data["role"], nil
}

func (s *SessionStore) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	key := fmt.Sprintf("refresh:%s", tokenHash)
	return s.client.Del(ctx, key).Err()
}

func (s *SessionStore) BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", jti)
	return s.client.Set(ctx, key, "1", ttl).Err()
}

func (s *SessionStore) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti)
	_, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}