package services

import (
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

type ModuleServiceRepository interface {
	SendTask(task *entities.TaskEntitie) error
	ReceiveTask() (*entities.TaskResultEntitie, error)
	IsAvailable() (bool, error)
	GetModuleById(id string) (valueobjects.CommandValue, error)
	GetAllModules() ([]valueobjects.CommandValue, error)
	GetModuleURL() (string, error)
}

