package validations

import (
	"errors"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

func ValidateTaskStatus(taskStatus string, newState valueobjects.TaskStatus) error {
	if taskStatus == "" {
		return errors.New("task status is required")
	}
	if !IsValidTaskStatus(taskStatus) {
		return errors.New("task status is invalid")
	}
	if taskStatus == string(newState) {
		return errors.New("task status is already " + string(newState))
	}
	return nil
}

func IsValidTaskStatus(taskStatus string) bool {
	for _, validTaskStatus := range valueobjects.ValidTaskStatus {
		if string(validTaskStatus) == taskStatus {
			return true
		}
	}
	return false
}