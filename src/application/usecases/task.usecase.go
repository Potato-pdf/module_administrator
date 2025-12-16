package usecases

import (
	"time"

	"luminaMO-orq-modAI/src/application/mappers"
	"luminaMO-orq-modAI/src/domain/DTO"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/repositories"
	orchestrator "luminaMO-orq-modAI/src/domain/services"
	validations "luminaMO-orq-modAI/src/domain/validations/command"
	"luminaMO-orq-modAI/src/domain/valueobjects"

	"github.com/google/uuid"
)

type SubmitTaskUseCase struct {
	taskRepo        repositories.TaskRepository
	orchestratorSvc orchestrator.OrchestratorService
}

func NewSubmitTaskUseCase(taskRepo repositories.TaskRepository, orchestratorSvc orchestrator.OrchestratorService) *SubmitTaskUseCase {
	return &SubmitTaskUseCase{
		taskRepo:        taskRepo,
		orchestratorSvc: orchestratorSvc,
	}
}

func (u *SubmitTaskUseCase) Execute(taskDTO DTO.TaskDTO) (*DTO.TaskDTO, error) {
	// 1. Validar comando
	if err := validations.ValidateCommand(taskDTO.Module); err != nil {
		return nil, err
	}

	// 2. Map DTO to Entity
	taskID := taskDTO.ID
	if taskID == "" {
		taskID = generateTaskID()
		taskDTO.ID = taskID
	}

	taskEntity := &entities.TaskEntitie{
		ID:         taskID,
		Module:     valueobjects.CommandValue(taskDTO.Module),
		CreatedAt:  time.Now(),
		Payload:    taskDTO.Payload,
		TaskStatus: valueobjects.Pending,
		RetryCount: 0,
		MaxRetries: taskDTO.MaxRetryCount,
	}

	// 3. Persistir tarea
	if err := u.taskRepo.CreateTask(taskEntity); err != nil {
		return nil, err
	}

	// 4. Enviar al orquestador (encolar en channel)
	if err := u.orchestratorSvc.SubmitTask(taskEntity); err != nil {
		return nil, err
	}

	// 5. Retornar DTO actualizado
	return mappers.MapTaskToDTO(taskEntity), nil
}

func generateTaskID() string {
	return uuid.New().String()
}
