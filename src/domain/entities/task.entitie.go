package entities

import (
	"time"

	command "luminaMO-orq-modAI/src/domain/valueobjects"
	taskStatus "luminaMO-orq-modAI/src/domain/valueobjects"
)

type TaskEntitie struct {
	ID          string                 
	Module     command.CommandValue   
	CreatedAt   time.Time              
	CompletedAt time.Time              
	Payload     map[string]interface{} 
	TaskStatus  taskStatus.TaskStatus 
	RetryCount  int                    
	MaxRetries  int                    
}

// =====Getters=====

func (t *TaskEntitie) GetID() string {
	return t.ID
}
	
func (t *TaskEntitie) GetModule() command.CommandValue {
	return t.Module
}

func (t *TaskEntitie) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *TaskEntitie) GetCompletedAt() time.Time {
	return t.CompletedAt
}

func (t *TaskEntitie) GetPayload() map[string]interface{} {
	return t.Payload
}	

func (t *TaskEntitie) GetTaskStatus() taskStatus.TaskStatus {
	return t.TaskStatus
}

func (t *TaskEntitie) GetRetryCount() int {
	return t.RetryCount
}

func (t *TaskEntitie) GetMaxRetries() int {
	return t.MaxRetries
}



	
