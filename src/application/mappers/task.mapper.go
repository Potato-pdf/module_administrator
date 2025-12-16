package mappers

import (
	"luminaMO-orq-modAI/src/domain/DTO"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

func MapTaskToDTO(task *entities.TaskEntitie) *DTO.TaskDTO {
	if task == nil {
		return nil
	}
	return &DTO.TaskDTO{
		ID:            task.ID,
		Module:        string(task.Module),
		CreatedAt:     task.CreatedAt,
		Payload:       task.Payload,
		TaskStatus:    task.TaskStatus,
		RetryCount:    task.RetryCount,
		MaxRetryCount: task.MaxRetries,
	}
}

func MapDTOToTask(dto *DTO.TaskDTO) *entities.TaskEntitie {
	if dto == nil {
		return nil
	}
	return &entities.TaskEntitie{
		ID:         dto.ID,
		Module:     valueobjects.CommandValue(dto.Module),
		CreatedAt:  dto.CreatedAt,
		Payload:    dto.Payload,
		TaskStatus: dto.TaskStatus,
		RetryCount: dto.RetryCount,
		MaxRetries: dto.MaxRetryCount,
	}
}

func MapTaskResultToDTO(taskResult *entities.TaskResultEntitie) *DTO.TaskResultDTO {
	if taskResult == nil {
		return nil
	}
	return &DTO.TaskResultDTO{
		ID:          taskResult.ID,
		TaskID:      taskResult.TaskID,
		Module:      string(taskResult.Module),
		CreatedAt:   taskResult.CreatedAt,
		CompletedAt: taskResult.CompletedAt,
		Payload:     taskResult.Payload,
		TaskStatus:  taskResult.TaskStatus,
		Error:       taskResult.Error,
	}
}

func MapDTOToTaskResult(dto *DTO.TaskResultDTO) *entities.TaskResultEntitie {
	if dto == nil {
		return nil
	}
	return &entities.TaskResultEntitie{
		ID:          dto.ID,
		TaskID:      dto.TaskID,
		Module:      valueobjects.CommandValue(dto.Module),
		CreatedAt:   dto.CreatedAt,
		CompletedAt: dto.CompletedAt,
		Payload:     dto.Payload,
		TaskStatus:  dto.TaskStatus,
		Error:       dto.Error,
	}
}
