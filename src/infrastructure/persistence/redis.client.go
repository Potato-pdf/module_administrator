package persistence

import (
	"context"
	"fmt"
	"time"

	"luminaMO-orq-modAI/src/domain/entities"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	redisEntitie *entities.RedisEntitie
	client       *redis.Client
}

func NewRedisClient(redisEntitie *entities.RedisEntitie) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisEntitie.Host + ":" + redisEntitie.Port,
		Password: redisEntitie.Password,
		DB:       redisEntitie.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &RedisClient{
		redisEntitie: redisEntitie,
		client:       client,
	}, nil
}

func CloseRedisClient(client *RedisClient) error {
	if client != nil {
		return client.client.Close()
	}
	return nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
