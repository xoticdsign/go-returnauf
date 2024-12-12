package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/xoticdsign/returnauf/models/responses"
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
func RunRedis(addr string, password string) (*Cache, error) {
	config := redis.Options{
		Addr:     addr,
		Password: password,
	}

	client := redis.NewClient(&config)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, redis.ErrClosed
	}
	return &Cache{cache: client}, nil
}

// Добавляет цитаты в тестовый Кэш
func (c *Cache) PopulateCache() {
	for _, quote := range responses.TestQuotes {
		c.cache.Set(context.Background(), strconv.Itoa(quote.ID), quote.Quote, time.Duration(time.Minute*5))
	}
}

// Уничтожает тестовый Кэш
func (c *Cache) TeardownCache() {
	client := c.cache
	client.FlushAll(context.Background())

	defer client.Close()
}

// Сохраняет данные в Кэш
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	err := c.cache.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return redis.ErrClosed
	}
	return nil
}

// Находит данные в Кэше
func (c *Cache) Get(key string) (string, error) {
	quote, err := c.cache.Get(context.Background(), key).Result()
	if err == redis.ErrClosed {
		return "", redis.ErrClosed
	}
	if err == redis.Nil {
		return "", redis.Nil
	}
	return quote, nil
}
