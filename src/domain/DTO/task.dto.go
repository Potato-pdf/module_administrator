package DTO

import (
	"time"

	taskStatus "luminaMO-orq-modAI/src/domain/valueobjects"
)

type TaskDTO struct {
	ID            string                 `json:"id"`
	Module        string                 `json:"module"`
	CreatedAt     time.Time              `json:"created_at"`
	Payload       map[string]interface{} `json:"payload"`
	TaskStatus    taskStatus.TaskStatus  `json:"status"`
	RetryCount    int                    `json:"retry_count"`
	MaxRetryCount int                    `json:"max_retry_count"`
}
