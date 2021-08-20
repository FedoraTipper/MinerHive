package storer

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(redisEndpoint, username, password string, selectedDB int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Username: username,
		Password: password,
		DB:       selectedDB,
	})

	err := redisTestConnection(client)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func redisTestConnection(client *redis.Client) error {
	ctx := context.Background()
	err := client.Set(ctx, "anthive_conn_test", "test", 5).Err()

	if err != nil {
		return err
	}

	return nil
}
