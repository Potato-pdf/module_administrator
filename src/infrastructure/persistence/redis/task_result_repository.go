package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/repositories"
)

type taskResultRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTaskResultRepository(client *redis.Client) repositories.TaskResultRepository {
	return &taskResultRepository{
		client: client,
		ttl:    1 * time.Hour,
	}
}

func (r *taskResultRepository) Save(result *entities.TaskResultEntitie) error {
	ctx := context.Background()

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal task result: %w", err)
	}

	key := fmt.Sprintf("result:%s", result.ID)
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return err
	}

	// También indexar por TaskID para búsqueda rápida
	taskKey := fmt.Sprintf("task_result:%s", result.TaskID)
	return r.client.Set(ctx, taskKey, result.ID, r.ttl).Err()
}

func (r *taskResultRepository) FindByID(id string) (*entities.TaskResultEntitie, error) {
	ctx := context.Background()

	key := fmt.Sprintf("result:%s", id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("task result not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task result: %w", err)
	}

	var result entities.TaskResultEntitie
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task result: %w", err)
	}

	return &result, nil
}

func (r *taskResultRepository) FindByTaskID(taskID string) (*entities.TaskResultEntitie, error) {
	ctx := context.Background()

	// Primero obtener el ID del resultado
	taskKey := fmt.Sprintf("task_result:%s", taskID)
	resultID, err := r.client.Get(ctx, taskKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("task result not found for task: %s", taskID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task result ID: %w", err)
	}

	// Luego obtener el resultado completo
	return r.FindByID(resultID)
}

func (r *taskResultRepository) Update(result *entities.TaskResultEntitie) error {
	return r.Save(result)
}

func (r *taskResultRepository) Delete(id string) error {
	ctx := context.Background()
	key := fmt.Sprintf("result:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *taskResultRepository) FindAll() ([]*entities.TaskResultEntitie, error) {
	ctx := context.Background()

	keys, err := r.client.Keys(ctx, "result:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get result keys: %w", err)
	}

	var results []*entities.TaskResultEntitie
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var result entities.TaskResultEntitie
		if err := json.Unmarshal(data, &result); err != nil {
			continue
		}

		results = append(results, &result)
	}

	return results, nil
}
