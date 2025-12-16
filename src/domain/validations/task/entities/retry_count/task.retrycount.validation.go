package validations

import (
	"errors"

	"luminaMO-orq-modAI/src/domain/entities"
)

func ValidateRetryCount(t *entities.TaskEntitie) error {
	if t.RetryCount < 0 {
		return errors.New("retry count is required")
	}
	if t.MaxRetries < 1 {
		return errors.New("max retries is required")
	}
	if t.MaxRetries > 10 {
		return errors.New("max retries cant be greater than 10")
	}
	if t.RetryCount > t.MaxRetries {
		return errors.New("retry count cant be greater than max retries")
	}
	return nil
}
