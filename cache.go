package jikan

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/minnasync/jikan-go/internal/redisx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type CacheOpts struct {
	TTL *time.Duration
}

// baseCache is a simple inferface for implementing methods for advanced caching.
type baseCache[T any] interface {
	// Get will get a value from the cache.
	Get(ctx context.Context, key string) (*T, error)
	// Set will set a value in the cache.
	Set(ctx context.Context, key string, value T, opts *CacheOpts) error
	// BulkSet will set multiple values in the cache.
	BulkSet(ctx context.Context, keyValues map[string]T, opts *CacheOpts) error
	// Delete will delete a value from the cache.
	Delete(ctx context.Context, key string) error
}

type AnimeCache interface {
	AnimeCache() baseCache[Anime]
	AnimeFullCache() baseCache[AnimeFull]

	GetAnime(ctx context.Context, id string) (*Anime, error)
	GetAnimeFull(ctx context.Context, id string) (*AnimeFull, error)
	SetAnime(ctx context.Context, data Anime) error
	SetAnimeFull(ctx context.Context, data AnimeFull) error
	BulkSetAnime(ctx context.Context, data []Anime) error
}

type animeCacheImpl struct {
	anime     baseCache[Anime]
	animeFull baseCache[AnimeFull]
}

func newAnimeCache(anime baseCache[Anime], animeFull baseCache[AnimeFull]) AnimeCache {
	return &animeCacheImpl{anime: anime, animeFull: animeFull}
}

func (c animeCacheImpl) AnimeCache() baseCache[Anime] {
	return c.anime
}

func (c animeCacheImpl) AnimeFullCache() baseCache[AnimeFull] {
	return c.animeFull
}

func (c animeCacheImpl) GetAnime(ctx context.Context, id string) (*Anime, error) {
	return c.anime.Get(ctx, "jikan:anime_"+id)
}

func (c animeCacheImpl) GetAnimeFull(ctx context.Context, id string) (*AnimeFull, error) {
	return c.animeFull.Get(ctx, "jikan:anime-full_"+id)
}

func (c animeCacheImpl) SetAnime(ctx context.Context, data Anime) error {
	return c.anime.Set(ctx, "jikan:anime_"+strconv.Itoa(data.MalID), data, &CacheOpts{
		TTL: new(time.Hour * 24),
	})
}

func (c animeCacheImpl) SetAnimeFull(ctx context.Context, data AnimeFull) error {
	// If the full value is fetched, we can set the base as well.
	// Makes sense to do since it's the same info.
	if err := c.SetAnime(ctx, data.Anime); err != nil {
		return err
	}

	return c.animeFull.Set(ctx, "jikan:anime-full_"+strconv.Itoa(data.MalID), data, &CacheOpts{
		TTL: new(time.Hour * 24),
	})
}

func (c animeCacheImpl) BulkSetAnime(ctx context.Context, data []Anime) error {
	entries := make(map[string]Anime, len(data))
	for _, entry := range data {
		entries["jikan:anime_"+strconv.Itoa(entry.MalID)] = entry
	}

	return c.anime.BulkSet(ctx, entries, nil)
}

type inMemoryCacheEntry[T any] struct {
	value   T
	expires time.Time
}

type inMemoryCacheImpl[T any] struct {
	mu      sync.RWMutex
	entries map[string]inMemoryCacheEntry[T]
}

func newInMemoryCache[T any]() baseCache[T] {
	return &inMemoryCacheImpl[T]{
		entries: make(map[string]inMemoryCacheEntry[T]),
	}
}

func (c *inMemoryCacheImpl[T]) Get(ctx context.Context, key string) (*T, error) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok {
		return nil, ErrCacheMiss
	}

	if !entry.expires.IsZero() && time.Now().After(entry.expires) {
		_ = c.Delete(ctx, key)
		return nil, ErrCacheMiss
	}

	return &entry.value, nil
}

func (c *inMemoryCacheImpl[T]) Set(ctx context.Context, key string, value T, opts *CacheOpts) error {
	entry := inMemoryCacheEntry[T]{
		value: value,
	}

	if opts != nil && opts.TTL != nil {
		entry.expires = time.Now().Add(*opts.TTL)
	}

	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()

	return nil
}

func (c *inMemoryCacheImpl[T]) BulkSet(ctx context.Context, keyValues map[string]T, opts *CacheOpts) error {
	var expiresAt time.Time
	if opts != nil && opts.TTL != nil {
		expiresAt = time.Now().Add(*opts.TTL)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range keyValues {
		entry := inMemoryCacheEntry[T]{
			value:   v,
			expires: expiresAt,
		}

		c.entries[k] = entry
	}

	return nil
}

func (c *inMemoryCacheImpl[T]) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)

	return nil
}

type DefaultCache struct {
	anime AnimeCache
}

// DefaultCache is a cache manager for an in-memory cache.
func NewCache() Caches {
	return &DefaultCache{
		anime: newAnimeCache(
			newInMemoryCache[Anime](),
			newInMemoryCache[AnimeFull](),
		),
	}
}

func (c *DefaultCache) Anime() AnimeCache {
	return c.anime
}

type redisJSONCacheImpl[T any] struct {
	sf     singleflight.Group
	client *redis.Client
}

// NewRedisCache will create a new cache manager for Redis that implements the Cache interface.
func newRedisJSONCache[T any](client *redis.Client) baseCache[T] {
	return &redisJSONCacheImpl[T]{client: client}
}

func (c *redisJSONCacheImpl[T]) Get(ctx context.Context, key string) (*T, error) {
	result, err, _ := c.sf.Do(key, func() (any, error) {
		var v T
		if err := redisx.JSONUnwrap(ctx, c.client, key, "$", &v); err != nil {
			return nil, err
		}

		return v, nil
	})

	if err != nil {
		return nil, err
	}

	value, ok := result.(T)
	if !ok {
		return nil, nil
	}

	return &value, nil
}

func (c *redisJSONCacheImpl[T]) Set(ctx context.Context, key string, value T, opts *CacheOpts) error {
	pipeline := c.client.Pipeline()

	pipeline.JSONSet(ctx, key, "$", value)

	if opts != nil && opts.TTL != nil {
		pipeline.Expire(ctx, key, *opts.TTL)
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (c *redisJSONCacheImpl[T]) BulkSet(ctx context.Context, keyValues map[string]T, opts *CacheOpts) error {
	pipeline := c.client.Pipeline()

	for key, value := range keyValues {
		pipeline.JSONSet(ctx, key, "$", value)

		if opts != nil && opts.TTL != nil {
			pipeline.Expire(ctx, key, *opts.TTL)
		}
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (c *redisJSONCacheImpl[T]) Delete(ctx context.Context, key string) error {
	return c.client.JSONDel(ctx, key, "$").Err()
}

type Caches interface {
	Anime() AnimeCache
}

type RedisJSONCache struct {
	anime AnimeCache
}

// RedisJSONCache is a cache manager for Redis.
//
// This will only use Redis' JSON commands, so using this will require your Redis
// instance to have the JSON module loaded. This is available on the redis-stack
// builds.
//
// When setting up your redis.conf, all you need to do is add this line:
// `loadmodule /opt/redis-stack/lib/rejson.so`
func NewRedisJSONCache(client *redis.Client) Caches {
	return &RedisJSONCache{
		anime: newAnimeCache(
			newRedisJSONCache[Anime](client),
			newRedisJSONCache[AnimeFull](client),
		),
	}
}

func (c *RedisJSONCache) Anime() AnimeCache {
	return c.anime
}
