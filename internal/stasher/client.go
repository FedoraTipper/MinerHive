package stasher

import (
	"context"
	"fmt"
	"time"

	"github.com/FedoraTipper/AntHive/pkg/models"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Stasher struct {
	redisClient *redis.Client
}

func (s *Stasher) NewRedisClient(host string, port uint, username, password string, selectedDB int) error {
	url := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Username: username,
		Password: password,
		DB:       selectedDB,
	})

	s.redisClient = client

	return s.redisTestConnection(client)
}

func (s *Stasher) redisTestConnection(client *redis.Client) error {
	ctx := context.Background()

	return client.Ping(ctx).Err()
}

func (s *Stasher) StashInterface(miner *models.MinerStats, expiration time.Duration) error {
	ctx := context.Background()

	zap.S().Infow("Stashing miner stats in RedisDB with no expiration", "Miner", miner.MinerName)
	err := s.redisClient.Set(ctx, miner.MinerName, miner, -1).Err()

	if err != nil {
		return err
	}

	return nil
}

func (s *Stasher) GetInterface(key string) (string, error) {
	ctx := context.Background()

	zap.S().Infof("Getting miner stats with key (%s) from RedisDB", key)
	i, err := s.redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		zap.S().Warnf("Key %s is missing from RedisDB", key)
		err = nil
	} else if err != nil {
		return "", err
	}

	return i, nil
}
