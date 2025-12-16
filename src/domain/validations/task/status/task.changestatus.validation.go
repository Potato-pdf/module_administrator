package validations

import (
	"errors"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

func ValidateFailTask(t *entities.TaskEntitie) error {
	if t.RetryCount >= t.MaxRetries {
		return errors.New("max retries reached")
	}
	if t.TaskStatus != valueobjects.Failed {
		return errors.New("task status is not failed")
	}
	if t.TaskStatus == valueobjects.Done {
		return errors.New("task status is done")
	}
	return nil
}


func CanTracitionTo(ts valueobjects.TaskStatus, newState valueobjects.TaskStatus) bool {
	trancition := map[valueobjects.TaskStatus][]valueobjects.TaskStatus{
	valueobjects.Pending:	{valueobjects.Processing, valueobjects.Failed},
	valueobjects.Processing:	{valueobjects.Done, valueobjects.Failed},
	valueobjects.Failed:	{valueobjects.Done},
	valueobjects.Done:	{},
	}
	allowedTransitions, ok := trancition[ts]
	if !ok {
		return false
	}
	for _, transition := range allowedTransitions {
		if transition == newState {
			return true
		}
	}
	return false
}


func ValidateChangeTaskStatus(t *entities.TaskEntitie, newState valueobjects.TaskStatus) error {
	if !CanTracitionTo(t.TaskStatus, newState) {
		return errors.New("task status transition is not allowed")
	}
	return nil
}

func IsTerminal(ts valueobjects.TaskStatus) bool {
	return ts == valueobjects.Done
}

func IsFailed(ts valueobjects.TaskStatus) bool {
	return ts == valueobjects.Failed
}

func IsPending(ts valueobjects.TaskStatus) bool {
	return ts == valueobjects.Pending
}
		