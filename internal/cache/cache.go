package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func RunRedis() error {
	config := redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	Cache = redis.NewClient(&config)

	err := Cache.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	return nil
}
