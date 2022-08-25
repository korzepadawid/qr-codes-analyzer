package cache

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/korzepadawid/qr-codes-analyzer/config"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(config *config.Config) *redisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
		DB:       0,
	})

	return &redisCache{
		client: rdb,
	}
}

func (c *redisCache) Get(key string) (string, error) {
	result, err := c.client.Get(context.Background(), key).Result()

	if err != nil {
		return "", ErrKeyNotFound
	}

	return result, nil
}

func (c *redisCache) Set(params *SetParams) error {
	return c.client.Set(context.Background(), params.Key, params.Value, params.Duration).Err()
}
