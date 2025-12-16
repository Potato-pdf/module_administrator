package repositories

import "luminaMO-orq-modAI/src/domain/entities"

type TaskRepository interface {
	CreateTask(task *entities.TaskEntitie) error
	FindTaskByID(id string) (*entities.TaskEntitie, error)
	FindAllTasks() ([]*entities.TaskEntitie, error)
	FindFailedTasks() ([]*entities.TaskEntitie, error)
	FindCompletedTasks() ([]*entities.TaskEntitie, error)
	UpdateTask(task *entities.TaskEntitie) error
	DeleteTask(id string) error
}
