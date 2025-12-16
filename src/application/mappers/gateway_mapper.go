package mappers

import (
	"luminaMO-orq-modAI/src/domain/DTO"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

// GatewayRequest representa el formato que envía el API Gateway Entrada
type GatewayRequest struct {
	TaskID        string                 `json:"taskId,omitempty"`
	Module        string                 `json:"module"`
	Payload       map[string]interface{} `json:"payload"`
	MaxRetryCount int                    `json:"maxRetryCount,omitempty"`
	Priority      string                 `json:"priority,omitempty"`
	CorrelationID string                 `json:"correlationId,omitempty"`
}

// GatewayResponse representa el formato que espera el API Gateway Salida
type GatewayResponse struct {
	TaskID        string                 `json:"taskId"`
	Status        string                 `json:"status"`
	Result        map[string]interface{} `json:"result,omitempty"`
	Error         string                 `json:"error,omitempty"`
	ProcessedAt   string                 `json:"processedAt"`
	CorrelationID string                 `json:"correlationId,omitempty"`
}

// MapGatewayRequestToTaskDTO convierte el request del Gateway a TaskDTO interno
func MapGatewayRequestToTaskDTO(req *GatewayRequest) *DTO.TaskDTO {
	if req == nil {
		return nil
	}

	maxRetries := req.MaxRetryCount
	if maxRetries == 0 {
		maxRetries = 3 // Default
	}

	return &DTO.TaskDTO{
		ID:            req.TaskID,
		Module:        req.Module,
		Payload:       req.Payload,
		MaxRetryCount: maxRetries,
		TaskStatus:    valueobjects.Pending,
	}
}

// MapTaskResultToGatewayResponse convierte TaskResultDTO a formato del Gateway Salida
func MapTaskResultToGatewayResponse(result *DTO.TaskResultDTO) *GatewayResponse {
	if result == nil {
		return nil
	}

	return &GatewayResponse{
		TaskID:      result.TaskID,
		Status:      string(result.TaskStatus),
		Result:      result.Payload,
		Error:       result.Error,
		ProcessedAt: result.CompletedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// MapTaskEntityToGatewayResponse convierte TaskEntity a formato del Gateway (para status)
func MapTaskEntityToGatewayResponse(task *entities.TaskEntitie) *GatewayResponse {
	if task == nil {
		return nil
	}

	return &GatewayResponse{
		TaskID:      task.ID,
		Status:      string(task.TaskStatus),
		ProcessedAt: task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
