package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() (*redis.Client, error) {
	url := "redis://redis:6379/0?protocol=3"

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis connect: %w", err)
	}

	rdb := redis.NewClient(opts)

	return rdb, nil
}
