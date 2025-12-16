package repositories

import "luminaMO-orq-modAI/src/domain/entities"

type TaskResultRepository interface {
	Save(result *entities.TaskResultEntitie) error
	FindByID(id string) (*entities.TaskResultEntitie, error)
	FindByTaskID(taskID string) (*entities.TaskResultEntitie, error)
	Update(result *entities.TaskResultEntitie) error
	Delete(id string) error
	FindAll() ([]*entities.TaskResultEntitie, error)
}
