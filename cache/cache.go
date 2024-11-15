package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func RunRedis() error {
	config := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	}

	Cache = redis.NewClient(&config)

	err := Cache.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	return nil
}
