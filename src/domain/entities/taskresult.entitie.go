package entities

import (
	"time"

	command "luminaMO-orq-modAI/src/domain/valueobjects"
	taskStatus "luminaMO-orq-modAI/src/domain/valueobjects"
)

type TaskResultEntitie struct {
	ID          string
	TaskID      string
	Module      command.CommandValue
	CreatedAt   time.Time
	CompletedAt time.Time
	Payload     map[string]interface{}
	TaskStatus  taskStatus.TaskStatus
	Error       string
}

func (t *TaskResultEntitie) GetID() string {
	return t.ID
}
	
func (t *TaskResultEntitie) GetTaskID() string {
	return t.TaskID
}

func (t *TaskResultEntitie) GetModule() command.CommandValue {
	return t.Module
}

func (t *TaskResultEntitie) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *TaskResultEntitie) GetCompletedAt() time.Time {
	return t.CompletedAt
}

func (t *TaskResultEntitie) GetPayload() map[string]interface{} {
	return t.Payload
}

func (t *TaskResultEntitie) GetTaskStatus() taskStatus.TaskStatus {
	return t.TaskStatus
}

func (t *TaskResultEntitie) GetError() string {
	return t.Error
}	