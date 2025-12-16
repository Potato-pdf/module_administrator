package orchestrator

import (
	"luminaMO-orq-modAI/src/domain/DTO"
)

type Channels struct {
	InboundTasks    chan *DTO.TaskDTO
	OutboundResults chan *DTO.TaskResultDTO
	FailedTasks     chan *DTO.TaskDTO
	Shutdown        chan struct{}
}

func NewChannels(bufferSize int) *Channels {
	return &Channels{
		InboundTasks:    make(chan *DTO.TaskDTO, bufferSize),
		OutboundResults: make(chan *DTO.TaskResultDTO, bufferSize),
		FailedTasks:     make(chan *DTO.TaskDTO, bufferSize),
		Shutdown:        make(chan struct{}),
	}
}

func (c *Channels) Close() {
	close(c.Shutdown)
	close(c.InboundTasks)
	close(c.OutboundResults)
	close(c.FailedTasks)
}
