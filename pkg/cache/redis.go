package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	"event-service/conf"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type RedisCache struct {
	client *redis.Client
}

func (r RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r RedisCache) Set(key string, value interface{}) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func NewRedisCache(cfg conf.Config) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return RedisCache{client: client}, nil
}
