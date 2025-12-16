package config

import (
	"log"

	"luminaMO-orq-modAI/src/application/usecases"
	"luminaMO-orq-modAI/src/domain/entities"
	"luminaMO-orq-modAI/src/domain/services"
	"luminaMO-orq-modAI/src/domain/valueobjects"
	"luminaMO-orq-modAI/src/infrastructure/clients"
	"luminaMO-orq-modAI/src/infrastructure/orchestrator"
	"luminaMO-orq-modAI/src/infrastructure/persistence"
	redisRepo "luminaMO-orq-modAI/src/infrastructure/persistence/redis"
	"luminaMO-orq-modAI/src/infrastructure/security"
	"luminaMO-orq-modAI/src/presentation/http/handlers"
)

type Dependencies struct {
	TaskHandler         *handlers.TaskHandler
	OrchestratorService services.OrchestratorService
}

func InitializeDependencies(cfg *Config) (*Dependencies, error) {
	// Redis Client
	redisClient, err := persistence.NewRedisClient(&entities.RedisEntitie{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return nil, err
	}
	log.Println("✅ Redis connected")

	// Repositories
	taskRepo := redisRepo.NewTaskRepository(redisClient.GetClient())
	taskResultRepo := redisRepo.NewTaskResultRepository(redisClient.GetClient())
	retryQueueRepo := redisRepo.NewRetryQueueRepository(redisClient.GetClient())

	// Auth Service
	authService := security.NewAuthService(security.AuthConfig{
		InboundKeys:  cfg.Security.InboundKeys,
		OutboundKeys: cfg.Security.OutboundKeys,
		Permissions: map[string][]string{
			cfg.Security.InboundKeys["gateway-entrada"]: {"*"},
			cfg.Security.InboundKeys["module-rag"]:      {"rag"},
			cfg.Security.InboundKeys["module-mcp"]:      {"mcp"},
		},
	})

	// Module URLs
	moduleURLs := map[valueobjects.CommandValue]string{
		valueobjects.CommandValue("RAG"): cfg.Modules.RAGURL,
		valueobjects.CommandValue("MCP"): cfg.Modules.MCPURL,
	}

	// HTTP Clients
	moduleClient := clients.NewModuleClient(authService, moduleURLs)
	gatewayClient := clients.NewGatewayClient(authService, cfg.Modules.GatewaySalidaURL)

	// Channels
	channels := orchestrator.NewChannels(100)

	// Orchestrator Service
	orchService := orchestrator.NewOrchestratorService(
		channels,
		taskRepo,
		taskResultRepo,
		retryQueueRepo,
		moduleClient,
		gatewayClient,
		cfg.Server.WorkerCount,
	)

	// Use Cases
	submitTaskUseCase := usecases.NewSubmitTaskUseCase(taskRepo, orchService)
	getTaskStatusUseCase := usecases.NewGetTaskStatusUseCase(taskRepo)

	// Handlers
	taskHandler := handlers.NewTaskHandler(submitTaskUseCase, getTaskStatusUseCase)

	return &Dependencies{
		TaskHandler:         taskHandler,
		OrchestratorService: orchService,
	}, nil
}
