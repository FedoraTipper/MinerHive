package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	redisClient *redis.Client
}

func (s *RedisClient) NewRedisClient(address string, username, password string, selectedDB int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       selectedDB,
	})

	s.redisClient = client

	return s.redisTestConnection(client)
}

func (s *RedisClient) redisTestConnection(client *redis.Client) error {
	ctx := context.Background()

	return client.Ping(ctx).Err()
}

func (s *RedisClient) StashInterface(key string, i interface{}, expiration time.Duration) error {
	return s.redisClient.Set(context.Background(), key, i, expiration).Err()
}

func (s *RedisClient) GetInterface(key string) (string, error) {
	return s.redisClient.Get(context.Background(), key).Result()
}

func (s *RedisClient) GetKeys() ([]string, error) {
	return s.redisClient.Keys(context.Background(), "*").Result()
}
