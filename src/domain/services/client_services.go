package services

import "luminaMO-orq-modAI/src/domain/entities"

type AuthService interface {
	ValidateAPIKey(apiKey string, expectedSource string) error
	GetAPIKeyForDestination(destination string) (string, error)
	ValidateModulePermission(apiKey string, module string) error
}

type ModuleClient interface {
	SendTask(task *entities.TaskEntitie) (*entities.TaskResultEntitie, error)
	IsAvailable() (bool, error)
}

type GatewayClient interface {
	SubmitResult(result *entities.TaskResultEntitie) error
	IsHealthy() (bool, error)
}
