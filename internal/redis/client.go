package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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
	ctx := context.Background()

	zap.S().Infow("Stashing interface in RedisDB with expiration", "Key", key, "Expiration", expiration.String())
	err := s.redisClient.Set(ctx, key, i, expiration).Err()

	if err != nil {
		return err
	}

	return nil
}

func (s *RedisClient) GetInterface(key string) (string, error) {
	ctx := context.Background()

	zap.S().Infof("Getting string interface with key (%s) from RedisDB", key)
	i, err := s.redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		zap.S().Warnf("Key %s is missing from RedisDB", key)
		err = nil
	} else if err != nil {
		return "", err
	}

	return i, nil
}

func (s *RedisClient) GetKeys() ([]string, error) {
	ctx := context.Background()

	zap.S().Info("Getting all keys from RedisDB")
	keys, err := s.redisClient.Keys(ctx, "*").Result()

	if err == redis.Nil {
		zap.S().Warn("No keys returned from RedisDB")
		err = nil
	} else if err != nil {
		return nil, err
	}

	return keys, nil
}
