package cache

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Интерфейс, содержащий методы для работы с Кэшом
type Cacher interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}

// Структура, реализующая Cacher
type Cache struct {
	cache *redis.Client
}

// Запускает Redis и возвращает структуру, реализующую Cacher
func RunRedis() (*Cache, error) {
	config := redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	redis := redis.NewClient(&config)

	err := redis.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &Cache{cache: redis}, nil
}

// Сохраняет данные в Кэш
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	status := c.cache.Set(context.Background(), key, value, expiration)
	if status.Err() != nil {
		return redis.TxFailedErr
	}
	return nil
}

// Находит данные в Кэше
func (c *Cache) Get(key string) (string, error) {
	quote, err := c.cache.Get(context.Background(), key).Result()
	if err == redis.TxFailedErr {
		return "", redis.TxFailedErr
	}
	if err == redis.Nil {
		return "", redis.Nil
	}
	return quote, nil
}
