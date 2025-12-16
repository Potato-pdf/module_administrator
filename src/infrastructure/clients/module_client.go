package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/services"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

type moduleClient struct {
	httpClient  *http.Client
	authService services.AuthService
	moduleURLs  map[valueobjects.CommandValue]string
}

func NewModuleClient(
	authService services.AuthService,
	moduleURLs map[valueobjects.CommandValue]string,
) services.ModuleClient {
	return &moduleClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authService: authService,
		moduleURLs:  moduleURLs,
	}
}

func (c *moduleClient) SendTask(task *entities.TaskEntitie) (*entities.TaskResultEntitie, error) {
	// 1. Obtener URL del módulo
	moduleURL, exists := c.moduleURLs[task.Module]
	if !exists {
		return nil, fmt.Errorf("unknown module: %s", task.Module)
	}

	// 2. Obtener API Key para este módulo
	destination := fmt.Sprintf("module-%s", task.Module)
	apiKey, err := c.authService.GetAPIKeyForDestination(destination)
	if err != nil {
		return nil, err
	}

	// 3. Preparar payload
	payload, err := json.Marshal(map[string]interface{}{
		"taskId":  task.ID,
		"payload": task.Payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 4. Crear request
	req, err := http.NewRequest("POST", moduleURL+"/process", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 5. Agregar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Task-ID", task.ID)
	req.Header.Set("X-Orchestrator-ID", "orquestador-main")

	// 6. Enviar request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("module returned status: %d", resp.StatusCode)
	}

	// 7. Parsear respuesta
	var result struct {
		TaskID string                 `json:"taskId"`
		Status string                 `json:"status"`
		Result map[string]interface{} `json:"result"`
		Error  string                 `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 8. Convertir a TaskResultEntity
	taskResult := &entities.TaskResultEntitie{
		ID:          fmt.Sprintf("result-%s", task.ID),
		TaskID:      task.ID,
		Module:      task.Module,
		Payload:     result.Result,
		TaskStatus:  valueobjects.TaskStatus(result.Status),
		Error:       result.Error,
		CompletedAt: time.Now(),
	}

	return taskResult, nil
}

func (c *moduleClient) IsAvailable() (bool, error) {
	// Implementar health check básico
	return true, nil
}
