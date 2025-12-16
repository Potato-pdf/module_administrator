package handlers

import (
	"net/http"

	"luminaMO-orq-modAI/src/application/usecases"
	"luminaMO-orq-modAI/src/domain/DTO"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	submitTaskUseCase    *usecases.SubmitTaskUseCase
	getTaskStatusUseCase *usecases.GetTaskStatusUseCase
}

func NewTaskHandler(
	submitTaskUseCase *usecases.SubmitTaskUseCase,
	getTaskStatusUseCase *usecases.GetTaskStatusUseCase,
) *TaskHandler {
	return &TaskHandler{
		submitTaskUseCase:    submitTaskUseCase,
		getTaskStatusUseCase: getTaskStatusUseCase,
	}
}

func (h *TaskHandler) SubmitTask(c *gin.Context) {
	var taskDTO DTO.TaskDTO

	if err := c.ShouldBindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.submitTaskUseCase.Execute(taskDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	result, err := h.getTaskStatusUseCase.Execute(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
