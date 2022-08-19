package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() *redisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
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
