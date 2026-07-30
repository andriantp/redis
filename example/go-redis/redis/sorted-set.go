package redis

import (
	"context"

	redisv9 "github.com/redis/go-redis/v9"
)

func (r *repository) ZADD(ctx context.Context, key string, score float64, member string) error {
	return r.client.ZAdd(
		ctx,
		key,
		redisv9.Z{
			Score:  score,
			Member: member,
		},
	).Err()
}

func (r *repository) ZREVRANGE(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRangeArgs(
		ctx,
		redisv9.ZRangeArgs{
			Key:   key,
			Start: start,
			Stop:  stop,
			Rev:   true,
		},
	).Result()
}
