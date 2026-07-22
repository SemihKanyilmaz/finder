package db

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(dsn string) (*redis.Client, error) {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("redis parse url: %w", err)
	}

	return redis.NewClient(opts), nil
}
