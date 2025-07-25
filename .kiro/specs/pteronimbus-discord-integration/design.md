# Design Document

## Overview

Pteronimbus is designed as a multi-tenant, Kubernetes-native game server management platform with deep Discord integration. The system follows a microservices architecture with clear separation of concerns between the web frontend, backend API, and Kubernetes controller. The platform leverages Discord's OAuth2 for authentication and uses Discord servers as natural tenant boundaries, creating isolated environments for each gaming community.

The architecture emphasizes security through isolation, scalability through Kubernetes-native patterns, and user experience through seamless Discord integration. The system is designed as a hosted SaaS platform with multi-tenant isolation.

## Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "User Interfaces"
        WEB[Web Frontend<br/>Nuxt.js]
        DISCORD[Discord Bot<br/>Commands & Notifications]
    end
    
    subgraph "Backend Services"
        API[Backend API<br/>Go/Gin]
        BOT[Discord Bot Service<br/>Go]
    end
    
    subgraph "Data Layer"
        DB[(PostgreSQL<br/>User Data & Config)]
        REDIS[(Redis<br/>Sessions & Cache)]
    end
    
    subgraph "Kubernetes Cluster"
        CONTROLLER[Pteronimbus Controller<br/>Go/Controller-Runtime]
        CRD1[GameTemplate CRDs]
        CRD2[GameServer CRDs]
        GAMESERVERS[Game Server Pods<br/>Minecraft, CS2, etc.]
    end
    
    subgraph "External Services"
        DISCORD_API[Discord API<br/>OAuth2 & Bot API]
    end
    
    WEB --> API
    DISCORD --> BOT
    BOT --> API
    API --> DB
    API --> REDIS
    API --> DISCORD_API
    CONTROLLER --> API
    CONTROLLER --> CRD1
    CONTROLLER --> CRD2
    CRD2 --> GAMESERVERS
    CRD1 -.-> CRD2
```

### Component Interaction Flow

1. **Authentication Flow**: Users authenticate via Discord OAuth2, backend issues JWT tokens
2. **Tenant Discovery**: Backend queries Discord API to determine user's accessible servers
3. **RBAC Sync**: System pulls Discord roles and creates internal permission mappings
4. **Game Server Management**: Users create GameServer CRDs via web UI or Discord commands
5. **Controller Reconciliation**: Controller watches CRDs and reconciles Kubernetes state
6. **Status Updates**: Real-time status flows from Kubernetes through controller to users

## Components and Interfaces

### Frontend Application (Nuxt.js)

**Purpose**: Web-based user interface for game server management

**Key Features**:
- Server-side rendering for optimal performance
- Responsive design for desktop and mobile
- Real-time updates via WebSocket connections
- Multi-tenant UI with tenant switching

**API Integration**:
- RESTful API calls to backend
- JWT token management with automatic refresh
- WebSocket connection for real-time updates

**Key Pages**:
- Login/OAuth callback handling
- Tenant selection dashboard
- Game server management interface
- RBAC configuration (for authorized users)
- Server logs and monitoring

### Backend API (Go/Gin)

**Purpose**: Central API server handling authentication, authorization, and business logic

**Key Modules**:

1. **Authentication Service**
   - Discord OAuth2 integration
   - JWT token generation and validation
   - Refresh token management
   - Session management with Redis

2. **Tenant Management Service**
   - Discord server discovery and validation
   - Tenant creation and configuration
   - User-tenant relationship management

3. **RBAC Service**
   - Discord role synchronization
   - Permission mapping and evaluation
   - Role-based access control enforcement

4. **Game Server Service**
   - GameTemplate and GameServer CRUD operations
   - Pterodactyl egg import and conversion
   - Server lifecycle management requests

5. **Discord Integration Service**
   - Discord API client management
   - Bot command processing
   - Webhook handling for role changes

**Database Schema**:
```sql
-- Core tenant and user management
tenants (id, discord_server_id, name, config, created_at, updated_at)
users (id, discord_user_id, username, avatar, created_at, updated_at)
user_tenants (user_id, tenant_id, roles, permissions, created_at)

-- Discord role mapping
discord_roles (id, tenant_id, discord_role_id, name, permissions)
discord_users (id, tenant_id, discord_user_id, roles, last_sync)

-- Game server management
game_templates (id, tenant_id, name, game_type, config, version, created_at)
game_servers (id, tenant_id, template_id, name, config, status, created_at)

-- Audit and sessions
audit_logs (id, user_id, tenant_id, action, resource, details, timestamp)
sessions (id, user_id, token_hash, expires_at, created_at)
```

### Discord Bot Service (Go)

**Purpose**: Handle Discord slash commands and provide notifications

**Key Features**:
- Slash command registration and handling
- Permission validation using backend RBAC
- Real-time status notifications
- Interactive command responses

**Commands**:
- `/server list` - List available game servers
- `/server start <name>` - Start a game server
- `/server stop <name>` - Stop a game server
- `/server status <name>` - Get server status
- `/server logs <name>` - Get recent server logs

### Kubernetes Controller (Go/Controller-Runtime)

**Purpose**: Reconcile desired state from backend API with actual Kubernetes resources

**Key Features**:
- Custom Resource Definition management
- Pull-based reconciliation with backend API
- Kubernetes resource lifecycle management
- Status reporting back to backend

**Reconciliation Logic**:
1. Watch for GameServer CRD changes
2. Query backend API for desired state
3. Compare with current Kubernetes state
4. Create/update/delete Kubernetes resources as needed
5. Update CRD status with current state
6. Report status back to backend API

**Managed Resources**:
- Deployments or StatefulSets for game servers
- Services for network access
- ConfigMaps for game configurations
- PersistentVolumeClaims for game data
- NetworkPolicies for tenant isolation

## Data Models

### Core Domain Models

```go
// Tenant represents a Discord server with Pteronimbus installed
type Tenant struct {
    ID              string    `json:"id"`
    DiscordServerID string    `json:"discord_server_id"`
    Name            string    `json:"name"`
    Config          TenantConfig `json:"config"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type TenantConfig struct {
    DefaultGameTemplate string            `json:"default_game_template"`
    ResourceLimits      ResourceLimits    `json:"resource_limits"`
    NotificationChannels []string         `json:"notification_channels"`
    Settings            map[string]string `json:"settings"`
}

// User represents a Discord user with access to tenants
type User struct {
    ID            string    `json:"id"`
    DiscordUserID string    `json:"discord_user_id"`
    Username      string    `json:"username"`
    Avatar        string    `json:"avatar"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// GameTemplate defines reusable game server configurations
type GameTemplate struct {
    ID          string                 `json:"id"`
    TenantID    string                 `json:"tenant_id"`
    Name        string                 `json:"name"`
    GameType    string                 `json:"game_type"`
    Version     string                 `json:"version"`
    Config      GameTemplateConfig     `json:"config"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
}

type GameTemplateConfig struct {
    Image           string            `json:"image"`
    Ports           []Port            `json:"ports"`
    Environment     map[string]string `json:"environment"`
    Resources       ResourceRequirements `json:"resources"`
    PersistentData  []VolumeMount     `json:"persistent_data"`
    StartupCommand  []string          `json:"startup_command"`
}

// GameServer represents an instance of a game server
type GameServer struct {
    ID         string            `json:"id"`
    TenantID   string            `json:"tenant_id"`
    TemplateID string            `json:"template_id"`
    Name       string            `json:"name"`
    Config     GameServerConfig  `json:"config"`
    Status     GameServerStatus  `json:"status"`
    CreatedAt  time.Time         `json:"created_at"`
    UpdatedAt  time.Time         `json:"updated_at"`
}

type GameServerStatus struct {
    Phase       string    `json:"phase"` // Pending, Running, Stopped, Failed
    Message     string    `json:"message"`
    LastUpdated time.Time `json:"last_updated"`
    PlayerCount int       `json:"player_count"`
    Uptime      duration  `json:"uptime"`
}
```

### Kubernetes Custom Resource Definitions

```yaml
# GameTemplate CRD
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: gametemplates.pteronimbus.io
spec:
  group: pteronimbus.io
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              gameType:
                type: string
              image:
                type: string
              ports:
                type: array
                items:
                  type: object
                  properties:
                    name:
                      type: string
                    port:
                      type: integer
                    protocol:
                      type: string
              resources:
                type: object
                properties:
                  requests:
                    type: object
                  limits:
                    type: object
          status:
            type: object
            properties:
              phase:
                type: string
              message:
                type: string

---
# GameServer CRD
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: gameservers.pteronimbus.io
spec:
  group: pteronimbus.io
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              templateRef:
                type: object
                properties:
                  name:
                    type: string
                  namespace:
                    type: string
              tenantId:
                type: string
              config:
                type: object
                additionalProperties: true
          status:
            type: object
            properties:
              phase:
                type: string
              message:
                type: string
              playerCount:
                type: integer
              uptime:
                type: string
```

## Error Handling

### API Error Responses

```go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// Standard error codes
const (
    ErrUnauthorized     = "UNAUTHORIZED"
    ErrForbidden        = "FORBIDDEN"
    ErrNotFound         = "NOT_FOUND"
    ErrValidation       = "VALIDATION_ERROR"
    ErrRateLimit        = "RATE_LIMIT_EXCEEDED"
    ErrDiscordAPI       = "DISCORD_API_ERROR"
    ErrKubernetesAPI    = "KUBERNETES_API_ERROR"
)
```

### Error Handling Strategies

1. **Authentication Errors**: Return 401 with clear message, trigger re-authentication flow
2. **Authorization Errors**: Return 403 with specific permission requirements
3. **Discord API Errors**: Implement exponential backoff, cache responses when possible
4. **Kubernetes API Errors**: Retry with backoff, report status to users
5. **Validation Errors**: Return detailed field-level error information
6. **Rate Limiting**: Implement per-user and per-tenant rate limits

### Circuit Breaker Pattern

Implement circuit breakers for external service calls:
- Discord API calls
- Kubernetes API calls
- Database connections

## Testing Strategy

### Unit Testing

**Backend Services**:
- Mock Discord API responses
- Test RBAC permission evaluation
- Test JWT token generation/validation
- Test database operations with test containers

**Controller**:
- Test reconciliation logic with fake Kubernetes client
- Test CRD status updates
- Test error handling and retry logic

**Frontend**:
- Component testing with Vue Test Utils
- API integration testing with MSW (Mock Service Worker)
- E2E testing with Playwright

### Integration Testing

**API Integration**:
- Test Discord OAuth2 flow end-to-end
- Test tenant creation and role synchronization
- Test game server lifecycle operations

**Controller Integration**:
- Test with real Kubernetes cluster (kind/minikube)
- Test CRD creation and reconciliation
- Test resource cleanup on deletion

### Load Testing

**Performance Targets**:
- Support 1000+ concurrent users per backend instance
- Handle 100+ game servers per tenant
- Discord role sync within 5 minutes
- API response times < 200ms (95th percentile)

**Load Testing Scenarios**:
- Concurrent user authentication
- Bulk game server operations
- Discord webhook processing
- Controller reconciliation under load

### Security Testing

**Authentication & Authorization**:
- Test JWT token validation and expiration
- Test RBAC permission enforcement
- Test tenant isolation
- Test Discord OAuth2 security

**Infrastructure Security**:
- Test Kubernetes RBAC configuration
- Test network policies for tenant isolation
- Test secrets management
- Test container security scanning

### Monitoring and Observability

**Metrics Collection**:
- Prometheus metrics for all services
- Custom metrics for game server status
- Discord API rate limit monitoring
- Kubernetes resource utilization

**Logging Strategy**:
- Structured logging with correlation IDs
- Centralized log aggregation
- Audit logging for all user actions
- Error tracking and alerting

**Health Checks**:
- Kubernetes liveness and readiness probes
- Discord API connectivity checks
- Database connection health
- Redis connectivity monitoring

**Alerting**:
- Service availability alerts
- High error rate alerts
- Discord API rate limit alerts
- Game server failure alerts