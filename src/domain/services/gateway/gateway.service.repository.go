package services

import "luminaMO-orq-modAI/src/domain/entities"

type GatewayServiceRepository interface {
	SubmitTask(task *entities.TaskResultEntitie) error
	IsHealthy() (bool, error)
}
