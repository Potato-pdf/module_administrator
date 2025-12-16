package channels

import (
	"luminaMO-orq-modAI/src/domain/DTO"
)

type OrchestatorChannelInterface struct {
	InTasks       chan *DTO.TaskDTO
	OutTasks      chan *DTO.TaskResultDTO
	FailedTasks   chan *DTO.TaskResultDTO
	FailedResults chan *DTO.TaskResultDTO
	Shutdown      chan struct{}
}
