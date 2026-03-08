package redis

import (
	"context"
	"fmt"
	"oolio/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService() (*RedisService, error) {
	env := config.Env()
	return &RedisService{
		client: redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", env.REDIS_HOST, env.REDIS_PORT),
			Password:    env.REDIS_PASSWORD,
			DB:          env.REDIS_DB,
			ReadTimeout: 25 * time.Second,
		}),
	}, nil
}

func (s *RedisService) Connect(ctx context.Context) error {
	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func (s *RedisService) Close() error {
	return s.client.Close()
}

func (s *RedisService) SetValue(ctx context.Context, key string, value interface{}) error {
	return s.SetValueWithExpiration(ctx, key, value, 0)
}

func (s *RedisService) SetValueWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := s.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}
	return nil
}

func (s *RedisService) GetValue(ctx context.Context, key string) (string, error) {
	value, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
	}
	return value, nil
}

func (s *RedisService) GetAll(ctx context.Context) ([]string, error) {
	var cursor uint64
	var result []string
	for {
		keys, err := s.client.Keys(ctx, "*").Result()
		if err == redis.Nil {
			return []string{}, fmt.Errorf("no entries found")
		}
		if err != nil {
			return []string{}, fmt.Errorf("failed to get keys %w", err)
		}

		for _, key := range keys {
			value, err := s.client.Get(ctx, key).Result()
			if err != nil {
				return []string{}, fmt.Errorf("failed to get value for key %s: %w", key, err)
			}
			result = append(result, value)
		}

		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func (s *RedisService) DeleteValue(ctx context.Context, key string) error {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}
