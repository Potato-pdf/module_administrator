### Prerrequisitos
- Go 1.21+
- Redis (via Docker)

### Instalación
 Iniciar Redis (en otro repositorio)
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
# Orchestator
