package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Redis    RedisConfig
	Security SecurityConfig
	Modules  ModulesConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	WorkerCount int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SecurityConfig struct {
	InboundKeys  map[string]string
	OutboundKeys map[string]string
}

type ModulesConfig struct {
	RAGURL           string
	MCPURL           string
	GatewaySalidaURL string
}

func LoadConfig() (*Config, error) {
	// Cargar .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	workerCount, _ := strconv.Atoi(getEnv("WORKER_COUNT", "5"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			WorkerCount: workerCount,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Security: SecurityConfig{
			InboundKeys: map[string]string{
				"gateway-entrada": getEnv("GATEWAY_ENTRADA_API_KEY", ""),
				"module-rag":      getEnv("MODULE_RAG_TO_ORQ_KEY", ""),
				"module-mcp":      getEnv("MODULE_MCP_TO_ORQ_KEY", ""),
			},
			OutboundKeys: map[string]string{
				"module-rag":     getEnv("ORQ_TO_MODULE_RAG_KEY", ""),
				"module-mcp":     getEnv("ORQ_TO_MODULE_MCP_KEY", ""),
				"gateway-salida": getEnv("ORQ_TO_GATEWAY_SALIDA_KEY", ""),
			},
		},
		Modules: ModulesConfig{
			RAGURL:           getEnv("MODULE_RAG_URL", "http://localhost:8081"),
			MCPURL:           getEnv("MODULE_MCP_URL", "http://localhost:8082"),
			GatewaySalidaURL: getEnv("GATEWAY_SALIDA_URL", "http://localhost:8083"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
