package DTO

import (
	"time"

	taskStatus "luminaMO-orq-modAI/src/domain/valueobjects"
)

type TaskResultDTO struct {
	ID          string                 `json:"id"`
	TaskID      string                 `json:"task_id"`
	Module      string                 `json:"module"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt time.Time              `json:"completed_at"`
	Payload     map[string]interface{} `json:"payload"`
	TaskStatus  taskStatus.TaskStatus  `json:"status"`
	Error       string                 `json:"error"`
}
