# jikan-go
A simple Jikan v4 client for Golang.

> [!CAUTION]
> This is not a production-ready implementation for use-cases outside of MinnaSync. Use this with immense caution. The library is immensely incomplete and lacks many features. If there is a feature that you need implemented, feel free to create a pull request implementing the feature. If you have a bug report for any existing features, please make an issue or pull request.

## Caching
By default, the library will do in-memory caching for results. At the moment, this cannot be disabled.
### Using Redis
The library supports Redis caching. You must have the JSON module loaded for this to work. A basic implementation is shown below.
```go
package main

import (
    "github.com/minnasync/jikan-go"
    "github.com/go-redis/redis/v9"
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
