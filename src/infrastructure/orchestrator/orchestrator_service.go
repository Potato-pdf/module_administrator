package orchestrator

import (
	"fmt"
	"log"
	"sync"
	"time"

	"luminaMO-orq-modAI/src/application/mappers"
	"luminaMO-orq-modAI/src/domain/DTO"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/repositories"
	"luminaMO-orq-modAI/src/domain/services"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

type orchestratorService struct {
	channels       *Channels
	taskRepo       repositories.TaskRepository
	taskResultRepo repositories.TaskResultRepository
	retryQueueRepo repositories.RetryQueueRepository
	moduleClient   services.ModuleClient
	gatewayClient  services.GatewayClient
	workerCount    int
	wg             sync.WaitGroup
	retryHandlerWg sync.WaitGroup
}

func NewOrchestratorService(
	channels *Channels,
	taskRepo repositories.TaskRepository,
	taskResultRepo repositories.TaskResultRepository,
	retryQueueRepo repositories.RetryQueueRepository,
	moduleClient services.ModuleClient,
	gatewayClient services.GatewayClient,
	workerCount int,
) services.OrchestratorService {
	return &orchestratorService{
		channels:       channels,
		taskRepo:       taskRepo,
		taskResultRepo: taskResultRepo,
		retryQueueRepo: retryQueueRepo,
		moduleClient:   moduleClient,
		gatewayClient:  gatewayClient,
		workerCount:    workerCount,
	}
}

func (o *orchestratorService) SubmitTask(task *entities.TaskEntitie) error {
	// Convertir a DTO y enviar al channel
	taskDTO := mappers.MapTaskToDTO(task)

	select {
	case o.channels.InboundTasks <- taskDTO:
		log.Printf("Task submitted to channel: %s", task.ID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout submitting task to channel")
	}
}

func (o *orchestratorService) Start() error {
	log.Printf("Starting orchestrator with %d workers", o.workerCount)

	// Iniciar workers
	for i := 0; i < o.workerCount; i++ {
		o.wg.Add(1)
		go o.worker(i)
	}

	// Iniciar retry handler
	o.retryHandlerWg.Add(1)
	go o.retryHandler()

	return nil
}

func (o *orchestratorService) Stop() error {
	log.Println("Stopping orchestrator...")
	o.channels.Close()
	o.wg.Wait()
	o.retryHandlerWg.Wait()
	log.Println("Orchestrator stopped")
	return nil
}

func (o *orchestratorService) worker(id int) {
	defer o.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case taskDTO := <-o.channels.InboundTasks:
			log.Printf("Worker %d processing task: %s", id, taskDTO.ID)
			o.processTask(taskDTO, id)

		case <-o.channels.Shutdown:
			log.Printf("Worker %d shutting down", id)
			return
		}
	}
}

func (o *orchestratorService) processTask(taskDTO *DTO.TaskDTO, workerID int) {
	// 1. Convertir DTO a Entity
	task := mappers.MapDTOToTask(taskDTO)

	// 2. Actualizar estado a Processing
	task.TaskStatus = valueobjects.Processing
	if err := o.taskRepo.UpdateTask(task); err != nil {
		log.Printf("Worker %d: Failed to update task status: %v", workerID, err)
	}

	// 3. Enviar tarea al módulo correspondiente
	result, err := o.moduleClient.SendTask(task)

	if err != nil {
		log.Printf("Worker %d: Module failed for task %s: %v", workerID, task.ID, err)
		o.handleTaskFailure(task, err, workerID)
		return
	}

	// 4. Procesar resultado exitoso
	log.Printf("Worker %d: Task %s completed successfully", workerID, task.ID)
	o.handleTaskSuccess(task, result, workerID)
}

func (o *orchestratorService) handleTaskSuccess(task *entities.TaskEntitie, result *entities.TaskResultEntitie, workerID int) {
	// 1. Actualizar estado de la tarea
	task.TaskStatus = valueobjects.Done
	if err := o.taskRepo.UpdateTask(task); err != nil {
		log.Printf("Worker %d: Failed to update task: %v", workerID, err)
	}

	// 2. Guardar resultado
	if err := o.taskResultRepo.Save(result); err != nil {
		log.Printf("Worker %d: Failed to save result: %v", workerID, err)
	}

	// 3. Enviar resultado al Gateway Salida
	if err := o.gatewayClient.SubmitResult(result); err != nil {
		log.Printf("Worker %d: Failed to submit result to gateway: %v", workerID, err)
		// Encolar para retry
		o.retryQueueRepo.EnqueueForRetry(task.ID, 30)
	}

	// 4. Enviar a channel de resultados
	resultDTO := mappers.MapTaskResultToDTO(result)
	select {
	case o.channels.OutboundResults <- resultDTO:
		log.Printf("Worker %d: Result sent to output channel", workerID)
	default:
		log.Printf("Worker %d: Output channel full, result not sent", workerID)
	}
}

func (o *orchestratorService) handleTaskFailure(task *entities.TaskEntitie, err error, workerID int) {
	// 1. Incrementar retry count
	task.RetryCount++

	// 2. Verificar si puede reintentar
	if task.RetryCount < task.MaxRetries {
		log.Printf("Worker %d: Task %s will retry (%d/%d)", workerID, task.ID, task.RetryCount, task.MaxRetries)

		// Actualizar estado a Failed (temporal)
		task.TaskStatus = valueobjects.Failed
		o.taskRepo.UpdateTask(task)

		// Encolar para retry (exponential backoff)
		retryDelay := 30 * task.RetryCount // 30s, 60s, 90s...
		o.retryQueueRepo.EnqueueForRetry(task.ID, retryDelay)

		// Enviar a channel de tareas fallidas
		taskDTO := mappers.MapTaskToDTO(task)
		select {
		case o.channels.FailedTasks <- taskDTO:
		default:
		}
	} else {
		log.Printf("Worker %d: Task %s exceeded max retries", workerID, task.ID)

		// Marcar como Failed permanentemente
		task.TaskStatus = valueobjects.Failed
		o.taskRepo.UpdateTask(task)

		// Crear resultado de error
		result := &entities.TaskResultEntitie{
			ID:          fmt.Sprintf("result-%s", task.ID),
			TaskID:      task.ID,
			Module:      task.Module,
			TaskStatus:  valueobjects.Failed,
			Error:       fmt.Sprintf("Max retries exceeded: %v", err),
			CompletedAt: time.Now(),
		}

		// Guardar resultado de error
		o.taskResultRepo.Save(result)

		// Notificar al Gateway Salida
		o.gatewayClient.SubmitResult(result)
	}
}

func (o *orchestratorService) retryHandler() {
	defer o.retryHandlerWg.Done()
	log.Println("Retry handler started")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.processRetries()

		case <-o.channels.Shutdown:
			log.Println("Retry handler shutting down")
			return
		}
	}
}

func (o *orchestratorService) processRetries() {
	// Obtener tareas listas para retry
	taskIDs, err := o.retryQueueRepo.GetTasksReadyForRetry(10)
	if err != nil {
		log.Printf("Retry handler: Failed to get tasks: %v", err)
		return
	}

	if len(taskIDs) == 0 {
		return
	}

	log.Printf("Retry handler: Processing %d tasks for retry", len(taskIDs))

	for _, taskID := range taskIDs {
		// Obtener tarea
		task, err := o.taskRepo.FindTaskByID(taskID)
		if err != nil {
			log.Printf("Retry handler: Task not found: %s", taskID)
			o.retryQueueRepo.RemoveFromRetryQueue(taskID)
			continue
		}

		// Remover de retry queue
		o.retryQueueRepo.RemoveFromRetryQueue(taskID)

		// Re-encolar en InboundTasks
		taskDTO := mappers.MapTaskToDTO(task)
		select {
		case o.channels.InboundTasks <- taskDTO:
			log.Printf("Retry handler: Task %s re-queued", taskID)
		default:
			log.Printf("Retry handler: Channel full, task %s not re-queued", taskID)
			// Volver a encolar para retry
			o.retryQueueRepo.EnqueueForRetry(taskID, 60)
		}
	}
}

func (o *orchestratorService) ProcessTask(taskID string) (*entities.TaskResultEntitie, error) {
	return o.taskResultRepo.FindByTaskID(taskID)
}

func (o *orchestratorService) GetTaskStatus(taskID string) (*entities.TaskEntitie, error) {
	return o.taskRepo.FindTaskByID(taskID)
}

func (o *orchestratorService) GetTaskResult(taskID string) (*entities.TaskResultEntitie, error) {
	return o.taskResultRepo.FindByTaskID(taskID)
}

func (o *orchestratorService) RetryFailedTask(taskID string) error {
	task, err := o.taskRepo.FindTaskByID(taskID)
	if err != nil {
		return err
	}

	// Reset retry count
	task.RetryCount = 0
	task.TaskStatus = valueobjects.Pending

	if err := o.taskRepo.UpdateTask(task); err != nil {
		return err
	}

	// Re-submit
	return o.SubmitTask(task)
}

func (o *orchestratorService) CancelTask(taskID string) error {
	task, err := o.taskRepo.FindTaskByID(taskID)
	if err != nil {
		return err
	}

	task.TaskStatus = valueobjects.Failed
	return o.taskRepo.UpdateTask(task)
}

func (o *orchestratorService) HealthCheck() error {
	// Verificar que los channels estén funcionando
	count, err := o.retryQueueRepo.CountPendingRetries()
	if err != nil {
		return fmt.Errorf("retry queue unhealthy: %w", err)
	}

	log.Printf("Health check: %d tasks in retry queue", count)
	return nil
}
