package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type Setting struct {
	Addr     string
	Username string
	Password string
	PathCert string
}

type repository struct {
	setting Setting
	client  *redisv9.Client
}

func NewRepository(setting Setting) (RepositoryI, error) {
	caCert, err := os.ReadFile(setting.PathCert)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	client := redisv9.NewClient(&redisv9.Options{
		Addr:     setting.Addr,
		Username: setting.Username,
		Password: setting.Password,

		TLSConfig: &tls.Config{
			RootCAs: caPool,
		},
	})
	return &repository{
		setting: setting,
		client:  client,
	}, nil
}

type RepositoryI interface {
	// Strings
	GET(ctx context.Context, key string) (string, error)
	SET(ctx context.Context, key, value string) error
	SETEX(ctx context.Context, key, value string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hashes
	HSET(ctx context.Context, key string, values map[string]interface{}) error
	HGET(ctx context.Context, key, field string) (string, error)
	HGETALL(ctx context.Context, key string) (map[string]string, error)

	// Pub/Sub
	PUBLISH(ctx context.Context, channel, message string) error
	SUBSCRIBE(ctx context.Context, channel string) error

	// Sorted
	ZADD(ctx context.Context, key string, score float64, member string) error
	ZREVRANGE(ctx context.Context, key string, start, stop int64) ([]string, error)
}
