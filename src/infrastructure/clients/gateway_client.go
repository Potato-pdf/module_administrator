package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/services"
)

type gatewayClient struct {
	httpClient  *http.Client
	authService services.AuthService
	gatewayURL  string
}

func NewGatewayClient(
	authService services.AuthService,
	gatewayURL string,
) services.GatewayClient {
	return &gatewayClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		authService: authService,
		gatewayURL:  gatewayURL,
	}
}

func (c *gatewayClient) SubmitResult(result *entities.TaskResultEntitie) error {
	// 1. Obtener API Key para Gateway Salida
	apiKey, err := c.authService.GetAPIKeyForDestination("gateway-salida")
	if err != nil {
		return err
	}

	// 2. Preparar payload
	payload, err := json.Marshal(map[string]interface{}{
		"taskId":      result.TaskID,
		"status":      string(result.TaskStatus),
		"result":      result.Payload,
		"error":       result.Error,
		"processedAt": result.CompletedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// 3. Crear request
	req, err := http.NewRequest("POST", c.gatewayURL+"/results", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Agregar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Task-ID", result.TaskID)
	req.Header.Set("X-Orchestrator-ID", "orquestador-main")

	// 5. Enviar request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("gateway returned status: %d", resp.StatusCode)
	}

	return nil
}

func (c *gatewayClient) IsHealthy() (bool, error) {
	// Implementar health check
	req, err := http.NewRequest("GET", c.gatewayURL+"/health", nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
