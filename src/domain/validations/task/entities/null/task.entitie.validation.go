package validations

import (
	"errors"

	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

func ValidateTaskEntitie(t *entities.TaskEntitie) error {
	if t.ID == "" {
		return errors.New("id is required")
	}
	if t.Module == "" {
		return errors.New("module is required")
	}
	if t.CreatedAt.IsZero() {
		return errors.New("created at is required")
	}
	if (t.TaskStatus == valueobjects.Done || t.TaskStatus == valueobjects.Failed) && t.CompletedAt.IsZero() {
		return errors.New("completed at is required for completed tasks")
	}
	if t.Payload == nil {
		return errors.New("payload is required")
	}
	if t.TaskStatus == "" {
		return errors.New("task status is required")
	}
	if t.RetryCount < 0 {
		return errors.New("retry count is required")
	}
	if t.MaxRetries < 0 {
		return errors.New("max retries is required")
	}
	return nil
}