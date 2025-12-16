package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"luminaMO-orq-modAI/src/domain/repositories"
)

type retryQueueRepository struct {
	client *redis.Client
}

func NewRetryQueueRepository(client *redis.Client) repositories.RetryQueueRepository {
	return &retryQueueRepository{
		client: client,
	}
}

func (r *retryQueueRepository) EnqueueForRetry(taskID string, retryAfterSeconds int) error {
	ctx := context.Background()

	// Calcular timestamp de retry
	retryAt := time.Now().Add(time.Duration(retryAfterSeconds) * time.Second).Unix()

	// Agregar a Sorted Set con score = timestamp
	return r.client.ZAdd(ctx, "retry_queue", &redis.Z{
		Score:  float64(retryAt),
		Member: taskID,
	}).Err()
}

func (r *retryQueueRepository) GetTasksReadyForRetry(limit int) ([]string, error) {
	ctx := context.Background()

	// Obtener tareas cuyo score (timestamp) <= ahora
	now := time.Now().Unix()

	taskIDs, err := r.client.ZRangeByScore(ctx, "retry_queue", &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(now, 10),
		Count: int64(limit),
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get tasks ready for retry: %w", err)
	}

	return taskIDs, nil
}

func (r *retryQueueRepository) RemoveFromRetryQueue(taskID string) error {
	ctx := context.Background()
	return r.client.ZRem(ctx, "retry_queue", taskID).Err()
}

func (r *retryQueueRepository) CountPendingRetries() (int64, error) {
	ctx := context.Background()
	return r.client.ZCard(ctx, "retry_queue").Result()
}

func (r *retryQueueRepository) IsInRetryQueue(taskID string) (bool, error) {
	ctx := context.Background()

	score, err := r.client.ZScore(ctx, "retry_queue", taskID).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return score > 0, nil
}
