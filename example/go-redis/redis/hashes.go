package redis

import (
	"context"
	"errors"
	"fmt"

	redisv9 "github.com/redis/go-redis/v9"
)

func (r *repository) HSET(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

func (r *repository) HGET(ctx context.Context, key, field string) (string, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redisv9.Nil) {
			return "", fmt.Errorf("field not found")
		}

		return "", err
	}

	return result, nil
}

func (r *repository) HGETALL(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}
