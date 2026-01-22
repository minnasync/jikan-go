# jikan-go
A simple Jikan v4 client for Golang.

> [!CAUTION]
> This is not a production-ready implementation for use-cases outside of MinnaSync. Use this with immense caution. The library is immensely incomplete and lacks many features. If there is a feature that you need implemented, feel free to create a pull request implementing the feature. If you have a bug report for any existing features, please make an issue or pull request.

## Caching
### Basic Caching
The library itself supports basic caching methods. At the moment, only Redis is supported. You must have the JSON module loaded for this to work. A basic implementation is shown below.
```go
package main

import (
    "github.com/MinnDevelopment/jikan-go"
)

func initRedis() *redis.Client {
	opts := &redis.Options{
        Addr:     "localhost:6379",
        DB:       0,
	}

	return redis.NewClient(opts)
}

func main() {
    redisClient := initRedis()
    client = jikan.NewClient(jikan.WithRedisCache(redisClient))
}
```

### Advanced Caching
If you have an advanced usecase and do not like how the caching works, you can implement your own caching manager by having it implement the `jikan.Cache` interface.
```go
package main

import (
    "github.com/MinnDevelopment/jikan-go"
)

type CustomCache struct {}

func NewCustomCache() jikan.Cache {
    return &CustomCache{}
}

func (c *CustomCache) Get(ctx context.Context, key string, v any) error {}
func (c *CustomCache) Set(ctx context.Context, key string, v any, ttl time.Duration) error {}
func (c *CustomCache) DeferSet(ctx context.Context, key string, v any, ttl time.Duration) {}
func (c *CustomCache) BulkSet(ctx context.Context, keyValues map[string]any, ttl time.Duration) error {}
func (c *CustomCache) DeferBulkSet(ctx context.Context, keyValues map[string]any, ttl time.Duration) {}
func (c *CustomCache) Delete(ctx context.Context, key string) error {}

func main() {
    cache := NewCustomCache()
    client = jikan.NewClient(jikan.WithCache(cache))
}
```