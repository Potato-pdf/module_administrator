package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"luminaMO-orq-modAI/src/config"
	"luminaMO-orq-modAI/src/presentation/http/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Inicializar dependencias
	deps, err := config.InitializeDependencies(cfg)
	if err != nil {
		log.Fatal("Failed to initialize dependencies:", err)
	}

	// Iniciar orchestrator
	if err := deps.OrchestratorService.Start(); err != nil {
		log.Fatal("Failed to start orchestrator:", err)
	}

	// Setup Gin
	router := gin.Default()
	routes.SetupRoutes(router, deps.TaskHandler)

	// HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	deps.OrchestratorService.Stop()
	log.Println("✅ Server exited")
}
