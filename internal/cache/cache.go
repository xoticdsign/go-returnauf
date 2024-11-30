package cache

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, string)
}

type Service struct {
	cache *redis.Client
}

func RunRedis() (*Service, error) {
	config := redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	redis := redis.NewClient(&config)

	err := redis.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &Service{cache: redis}, nil
}

func (s *Service) Set(key string, value interface{}, expiration time.Duration) error {
	status := s.cache.Set(context.Background(), key, value, expiration)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (s *Service) Get(key string) (string, string) {
	quote, err := s.cache.Get(context.Background(), key).Result()
	if err == redis.TxFailedErr {
		return "", "Failed"
	}
	if err == redis.Nil {
		return "", "Nil"
	}
	return quote, ""
}
