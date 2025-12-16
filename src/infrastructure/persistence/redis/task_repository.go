package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/repositories"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

type taskRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTaskRepository(client *redis.Client) repositories.TaskRepository {
	return &taskRepository{
		client: client,
		ttl:    1 * time.Hour,
	}
}

func (r *taskRepository) CreateTask(task *entities.TaskEntitie) error {
	ctx := context.Background()

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	key := fmt.Sprintf("task:%s", task.ID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *taskRepository) FindTaskByID(id string) (*entities.TaskEntitie, error) {
	ctx := context.Background()

	key := fmt.Sprintf("task:%s", id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var task entities.TaskEntitie
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) UpdateTask(task *entities.TaskEntitie) error {
	return r.CreateTask(task)
}

func (r *taskRepository) FindTasksByStatus(status valueobjects.TaskStatus) ([]*entities.TaskEntitie, error) {
	ctx := context.Background()

	keys, err := r.client.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task keys: %w", err)
	}

	var tasks []*entities.TaskEntitie
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var task entities.TaskEntitie
		if err := json.Unmarshal(data, &task); err != nil {
			continue
		}

		if task.TaskStatus == status {
			tasks = append(tasks, &task)
		}
	}

	return tasks, nil
}

func (r *taskRepository) DeleteTask(id string) error {
	ctx := context.Background()
	key := fmt.Sprintf("task:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *taskRepository) FindAllTasks() ([]*entities.TaskEntitie, error) {
	ctx := context.Background()

	keys, err := r.client.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task keys: %w", err)
	}

	var tasks []*entities.TaskEntitie
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var task entities.TaskEntitie
		if err := json.Unmarshal(data, &task); err != nil {
			continue
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) FindFailedTasks() ([]*entities.TaskEntitie, error) {
	return r.FindTasksByStatus(valueobjects.Failed)
}

func (r *taskRepository) FindCompletedTasks() ([]*entities.TaskEntitie, error) {
	return r.FindTasksByStatus(valueobjects.Done)
}
