package jikan

import (
	"context"
	"time"

	"github.com/minnasync/jikan-go/internal/redisx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Cache interface {
	// Get will get a value from the cache.
	Get(ctx context.Context, key string, value any) error
	// Set will set a value in the cache.
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	// Delete will delete a value from the cache.
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	sf     singleflight.Group
	client *redis.Client
}

// NewRedisCache will create a new cache manager for Redis that implements the Cache interface.
//
// This will only use JSON commands, so your Redis server must have the JSON module loaded.
func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string, v any) error {
	_, err, _ := c.sf.Do(key, func() (any, error) {
		data := new(any)

		if err := redisx.JSONUnwrap(ctx, c.client, key, "$", &data); err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return data, nil
	})

	return err
}

func (c *RedisCache) Set(ctx context.Context, key string, v any, ttl time.Duration) error {
	pipeline := c.client.Pipeline()

	pipeline.JSONSet(ctx, key, "$", v)
	pipeline.Expire(ctx, key, ttl)

	_, err := pipeline.Exec(ctx)
	return err
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	_, err, _ := c.sf.Do(key, func() (any, error) {
		cmd := c.client.JSONDel(ctx, key, "$")

		if err := cmd.Err(); err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil
	})

	return err
}
