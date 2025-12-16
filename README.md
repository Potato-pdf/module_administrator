# Lumina Orchestrator

Orquestador de tareas para módulos de IA (RAG, MCP, etc.)

## 🚀 Inicio Rápido

### Prerrequisitos
- Go 1.21+
- Redis (via Docker)

### Instalación

1. Clonar repositorio
```bash
git clone <repo-url>
cd LuminaMO_Orchestrator_Modules_Agents
```

2. Instalar dependencias
```bash
go mod download
```

3. Configurar variables de entorno
```bash
cp .env.example .env
# Editar .env con tus valores
```

4. Iniciar Redis (en otro repositorio)
```bash
docker-compose up -d
```

5. Ejecutar orquestador
```bash
go run main.go
```

## 📡 API Endpoints

### Health Check
```bash
GET /health
```

### Submit Task
```bash
POST /api/v1/tasks
Content-Type: application/json
X-API-Key: <tu-api-key>

{
  "module": "RAG",
  "payload": {
    "query": "¿Qué es Clean Architecture?"
  },
  "maxRetryCount": 3
}
```

### Get Task Status
```bash
GET /api/v1/tasks/:id
X-API-Key: <tu-api-key>
```

## 🏗️ Arquitectura

```
Clean Architecture:
- Domain: Entities, Value Objects, Interfaces
- Application: Use Cases, Mappers
- Infrastructure: Redis, HTTP Clients, Orchestrator
- Presentation: HTTP Handlers, Routes
```

## 📝 Licencia

MIT
# Orchestator
