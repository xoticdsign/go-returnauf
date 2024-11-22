package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func RunRedis() error {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	config := redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	}

	Cache = redis.NewClient(&config)

	err := Cache.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	return nil
}
