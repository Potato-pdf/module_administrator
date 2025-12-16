package services

import "luminaMO-orq-modAI/src/domain/entities"

type OrchestratorService interface {
	SubmitTask(task *entities.TaskEntitie) error
	ProcessTask(taskID string) (*entities.TaskResultEntitie, error)
	GetTaskStatus(taskID string) (*entities.TaskEntitie, error)
	GetTaskResult(taskID string) (*entities.TaskResultEntitie, error)
	RetryFailedTask(taskID string) error
	CancelTask(taskID string) error
	Start() error
	Stop() error
	HealthCheck() error
}
