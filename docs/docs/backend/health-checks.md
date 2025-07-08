# Backend Health Checks

The Pteronimbus backend service implements comprehensive health check endpoints for monitoring and Kubernetes integration.

## Implementation

The health check service is implemented in Go using the Gorilla Mux router and provides multiple endpoints for different monitoring scenarios.

### Health Check Endpoints

| Endpoint | Purpose | Response |
|----------|---------|----------|
| `/health` | General health status | Full health information with version |
| `/healthz` | Kubernetes-style health check | Same as `/health` |
| `/ready` | Readiness probe | Indicates service ready for traffic |
| `/live` | Liveness probe | Indicates service is alive |

### Response Format

All health check endpoints return JSON responses with consistent structure:

```json
{
  "status": "healthy",
  "timestamp": "2025-07-08T19:44:27.814Z",
  "service": "pteronimbus-backend",
  "version": "0.1.0"
}
```

### Status Values

- **healthy**: Service is fully operational (general health)
- **ready**: Service is ready to receive traffic (readiness probe)
- **alive**: Service process is running (liveness probe)

## Usage Examples

### Testing Health Endpoints

```bash
# General health check
curl http://localhost:8080/health

# Kubernetes-style health check  
curl http://localhost:8080/healthz

# Readiness probe
curl http://localhost:8080/ready

# Liveness probe  
curl http://localhost:8080/live
```

### Expected Responses

**General Health** (`/health`, `/healthz`):
```json
{
  "status": "healthy",
  "timestamp": "2025-07-08T19:44:27.814Z",
  "service": "pteronimbus-backend",
  "version": "0.1.0"
}
```

**Readiness Probe** (`/ready`):
```json
{
  "status": "ready", 
  "timestamp": "2025-07-08T19:44:27.814Z",
  "service": "pteronimbus-backend"
}
```

**Liveness Probe** (`/live`):
```json
{
  "status": "alive",
  "timestamp": "2025-07-08T19:44:27.814Z", 
  "service": "pteronimbus-backend"
}
```

## Kubernetes Integration

These endpoints are designed for Kubernetes health checks:

```yaml
# Example Kubernetes deployment health checks
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pteronimbus-backend
spec:
  template:
    spec:
      containers:
      - name: backend
        image: pteronimbus/backend:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready  
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Development

### Running the Service

```bash
cd apps/backend
go run cmd/server/main.go
```

The service starts on port 8080 and immediately responds to health checks.

### Docker Support

The service includes a Dockerfile with built-in health checks:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

## Implementation Details

The health check handlers are implemented in `internal/handlers/health.go` with:

- **Graceful shutdown**: Server handles SIGINT/SIGTERM properly
- **Structured responses**: Consistent JSON format across all endpoints  
- **HTTP status codes**: Proper 200 OK responses for healthy state
- **Timestamps**: UTC timestamps for monitoring and debugging
- **Service identification**: Clear service name and version info

This foundation provides the monitoring and health check capabilities needed for the full Pteronimbus system as development continues. 