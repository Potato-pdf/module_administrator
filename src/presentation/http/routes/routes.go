package routes

import (
	"luminaMO-orq-modAI/src/presentation/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	// Health check
	router.GET("/health", taskHandler.HealthCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Tasks
		v1.POST("/tasks", taskHandler.SubmitTask)
		v1.GET("/tasks/:id", taskHandler.GetTaskStatus)
	}
}
