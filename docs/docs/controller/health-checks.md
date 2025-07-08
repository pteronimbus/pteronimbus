# Controller Health Checks

The Pteronimbus Controller provides comprehensive health check endpoints for monitoring and observability in Kubernetes environments.

## Health Endpoints

The controller exposes the following health check endpoints on port **8080**:

| Endpoint | Purpose | Use Case |
|----------|---------|----------|
| `/health` | General health status | Load balancer health checks |
| `/healthz` | Kubernetes-style health | Standard k8s health endpoint |
| `/ready` | Readiness probe | Service ready to handle requests |
| `/live` | Liveness probe | Process is alive and responsive |

## Endpoint Details

### `/health` and `/healthz`

General health check endpoints that return the controller's operational status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-07-08T20:29:56.398126854Z",
  "service": "pteronimbus-controller", 
  "version": "0.1.0"
}
```

**HTTP Status:** `200 OK`

### `/ready`

Readiness probe endpoint indicating the controller is ready to process requests and watch for CRD changes.

**Response:**
```json
{
  "status": "ready",
  "timestamp": "2025-07-08T20:30:24.913767180Z",
  "service": "pteronimbus-controller"
}
```

**HTTP Status:** `200 OK`

### `/live`

Liveness probe endpoint confirming the controller process is alive and responding.

**Response:**
```json
{
  "status": "alive", 
  "timestamp": "2025-07-08T20:30:28.619967421Z",
  "service": "pteronimbus-controller"
}
```

**HTTP Status:** `200 OK`

## Kubernetes Integration

### Deployment Health Checks

When deploying the controller in Kubernetes, configure health checks in your deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pteronimbus-controller
spec:
  template:
    spec:
      containers:
      - name: controller
        image: pteronimbus/controller:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### Service Monitor

For Prometheus monitoring integration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: pteronimbus-controller
spec:
  selector:
    matchLabels:
      app: pteronimbus-controller
  endpoints:
  - port: http
    path: /health
    interval: 30s
```

## Docker Health Checks

The controller's Dockerfile includes a built-in health check:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

This ensures container orchestrators can detect unhealthy controller instances.

## Monitoring Recommendations

### Health Check Frequency

- **Liveness**: Every 30 seconds (detect hung processes)
- **Readiness**: Every 10 seconds (quick traffic routing)
- **External monitoring**: Every 60 seconds (reduce noise)

### Alert Thresholds

- **Critical**: Health check failures > 3 consecutive attempts
- **Warning**: Response time > 1 second
- **Info**: Controller restarts or ready state changes

### Observability

The controller logs health check requests at debug level. Enable with:

```bash
# Set log level to debug
export LOG_LEVEL=debug
```

Future enhancements will include:
- Metrics endpoint (`/metrics`) for Prometheus
- Detailed health status including CRD connectivity
- Custom health checks for controller-specific operations

## Testing

Health checks can be tested locally:

```bash
# Start the controller
go run ./cmd/server

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready  
curl http://localhost:8080/live
curl http://localhost:8080/healthz
```

Unit tests verify all health endpoints return proper JSON responses and HTTP status codes. 