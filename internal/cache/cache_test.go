package cache

/*

ДЛЯ ИСПОЛЬЗОВАНИЯ ДАННЫХ ТЕСТОВ НЕОБХОДИМО ИМЕТЬ
РАБОТАЮЩИЙ ЛОКАЛЬНЫЙ REDIS С АДРЕСОМ 127.0.0.1:6379
И БЕЗ ПАРОЛЯ !!!

*/

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Настройка Redis для тестов
func setupTestCache(emptyCache bool) *Cache {
	Cache, _ := RunRedis("127.0.0.1:6379", "")

	if !emptyCache {
		Cache.PopulateCache()
	}

	return Cache
}

// Unit тест для функции RunRedis
func TestUnitRunRedis(t *testing.T) {
	cases := []struct {
		name                    string
		addr                    string
		password                string
		wantRunRedisToReturnErr error
	}{
		{
			name:                    "general case",
			addr:                    "127.0.0.1:6379",
			password:                "",
			wantRunRedisToReturnErr: nil,
		},
		{
			name:                    "wrong address case",
			addr:                    "wrongaddr",
			password:                "",
			wantRunRedisToReturnErr: redis.ErrClosed,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			gotCache, gotErr := RunRedis(cs.addr, cs.password)
			if gotErr != nil {
				assert.Equalf(t, cs.wantRunRedisToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantRunRedisToReturnErr)
			} else {
				defer gotCache.TeardownCache()
			}
		})
	}
}

// Unit тест для функции Set
func TestUnitSet(t *testing.T) {
	cases := []struct {
		name               string
		key                string
		value              string
		wantSetToReturnErr error
	}{
		{
			name:               "general case",
			key:                "key",
			value:              "value",
			wantSetToReturnErr: nil,
		},
		{
			name:               "closed client case",
			key:                "key",
			value:              "value",
			wantSetToReturnErr: redis.ErrClosed,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			Cache := setupTestCache(false)
			client := Cache.cache
			defer Cache.TeardownCache()

			if cs.wantSetToReturnErr == redis.ErrClosed {
				client.Close()
			}

			gotErr := Cache.Set(cs.key, cs.value, time.Duration(time.Minute*5))
			if gotErr != nil {
				assert.Equalf(t, cs.wantSetToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantSetToReturnErr)
			} else {
				gotValue := Cache.cache.Get(context.Background(), cs.key).Val()

				assert.Equalf(t, cs.value, gotValue, "got %v, while comparing recently set value, want %v", gotValue, cs.value)
			}
		})
	}
}

// Unit тест для функции Get
func TestUnitGet(t *testing.T) {
	cases := []struct {
		name                 string
		key                  string
		emptyCache           bool
		wantGetToReturnValue string
		wantGetToReturnErr   error
	}{
		{
			name:                 "general case",
			key:                  "1",
			emptyCache:           false,
			wantGetToReturnValue: "Mock quote 1",
			wantGetToReturnErr:   nil,
		},
		{
			name:                 "empty cache case",
			key:                  "1",
			emptyCache:           true,
			wantGetToReturnValue: "",
			wantGetToReturnErr:   redis.Nil,
		},
		{
			name:                 "closed client case",
			key:                  "1",
			emptyCache:           false,
			wantGetToReturnValue: "",
			wantGetToReturnErr:   redis.ErrClosed,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			Cache := setupTestCache(cs.emptyCache)
			client := Cache.cache
			defer Cache.TeardownCache()

			if cs.wantGetToReturnErr == redis.ErrClosed {
				client.Close()
			}

			gotValue, gotErr := Cache.Get(cs.key)
			if gotErr != nil {
				assert.Equalf(t, cs.wantGetToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantGetToReturnErr)
			} else {
				assert.Equalf(t, cs.wantGetToReturnValue, gotValue, "got %v, while comparing returned value, want %v", gotValue, cs.wantGetToReturnValue)
			}
		})
	}
}
