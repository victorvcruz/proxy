package redis

import (
	"github.com/go-redis/redis/v8"
	"log"
	"proxy_project/cache"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) ConnectToDatabase() error {

	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis",
		DB:       0,
	})

	_, err := r.client.Ping(r.client.Context()).Result()
	if err != nil {
		return err
	}

	log.Println("Redis Connected")

	return nil
}

func (r *RedisClient) InsertInDatabase(key string, value string) error {
	err := r.client.Set(r.client.Context(), key, value, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) FindInDatabase(key string) (string, error) {
	val, err := r.client.Get(r.client.Context(), key).Result()
	if err != nil {
		return "", &cache.CacheNotFoundError{}
	}
	return val, nil
}
