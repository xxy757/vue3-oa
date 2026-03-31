package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(host string, port int, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Get(key string) (interface{}, bool) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}
