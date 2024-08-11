package cache

import (
	"context"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is a struct that represents a Redis cache.
type RedisCache struct {
	client *redis.Client
	// The default expiration time.
	defaultExpiration time.Duration
}

// NewRedisCache creates a new Redis cache with the given URL.
func NewRedisCache(redisConnStr string, defaultExpiration time.Duration) *RedisCache {
	redisURL, err := url.Parse(redisConnStr)
	if err != nil {
		panic(err)
	}

	rPassword, _ := redisURL.User.Password()
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Username: redisURL.User.Username(),
		Password: rPassword,
		DB:       0,
	})
	return &RedisCache{client: client, defaultExpiration: defaultExpiration}
}

// Get a value from the cache.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set a value in the cache.
func (r *RedisCache) Set(ctx context.Context, key, value string) error {
	err := r.client.Set(ctx, key, value, r.defaultExpiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete a value from the cache.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
