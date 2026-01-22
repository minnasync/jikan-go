# jikan-go
A simple Jikan v4 client for Golang.

> [!CAUTION]
> This is not a production-ready implementation for use-cases outside of MinnaSync. Use this with immense caution. The library is immensely incomplete and lacks many features. If there is a feature that you need implemented, feel free to create a pull request implementing the feature. If you have a bug report for any existing features, please make an issue or pull request.

## Features
- Ability to cache API responses using Redis. Client must be initalized with `WithRedisCache` to enable caching. Must be using `github.com/redis/go-redis/v9`. In the future, in-memory caching will also be supported.
```go
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
client := jikan.NewJikanClient(jikan.WithRedisCache(redisClient))
```