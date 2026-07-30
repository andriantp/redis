package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

func (r *repository) GET(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redisv9.Nil) {
			return "", fmt.Errorf("key not found")
		}

		return "", err
	}

	return result, nil
}

func (r *repository) SET(ctx context.Context, key, value string) error {
	return r.client.Set(
		ctx,
		key,
		value,
		0).
		Err()
}

func (r *repository) SETEX(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}

func (r *repository) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return ttl, nil
}
