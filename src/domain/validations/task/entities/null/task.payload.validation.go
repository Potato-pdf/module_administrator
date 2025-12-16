package validations

import (
	"errors"

	"luminaMO-orq-modAI/src/domain/entities"
)

func ValidatePayload(t *entities.TaskEntitie) error {
	if t.Payload == nil {
		return errors.New("payload is required")
	}
	if len(t.Payload) == 0 {
		return errors.New("payload is required")
	}
	if _, exists := t.Payload["action"]; !exists {
		return errors.New("action is required")
	}
	return nil
}
