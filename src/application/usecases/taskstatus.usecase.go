package usecases

import (
	"errors"
	"luminaMO-orq-modAI/src/application/mappers"
	"luminaMO-orq-modAI/src/domain/DTO"
	"luminaMO-orq-modAI/src/domain/repositories"
)

type GetTaskStatusUseCase struct {
	taskRepo repositories.TaskRepository
}

func NewGetTaskStatusUseCase(taskRepo repositories.TaskRepository) *GetTaskStatusUseCase {
	return &GetTaskStatusUseCase{
		taskRepo: taskRepo,
	}
}

func (u *GetTaskStatusUseCase) Execute(taskID string) (*DTO.TaskDTO, error) {
	if taskID == "" {
		return nil, errors.New("taskID is required")
	}

	task, err := u.taskRepo.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	
	return mappers.MapTaskToDTO(task), nil
}
