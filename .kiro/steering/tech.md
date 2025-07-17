# Technology Stack

## Architecture
- **Microservices**: Multi-app monorepo with separate frontend, backend, and controller services
- **Containerization**: Docker containers with Docker Compose for local development
- **Orchestration**: Kubernetes-native with Custom Resource Definitions (CRDs)
- **Deployment**: Helm charts for Kubernetes deployment

## Frontend Stack
- **Framework**: Nuxt 3 (Vue.js-based)
- **UI Library**: Nuxt UI with Tailwind CSS
- **Styling**: Tailwind CSS v4 with forms and typography plugins
- **Icons**: Nuxt Icon with Heroicons and Lucide
- **Internationalization**: Nuxt i18n
- **Testing**: Vitest with Vue Test Utils and Happy DOM
- **TypeScript**: Full TypeScript support

## Backend Stack
- **Language**: Go 1.23
- **Web Framework**: Gin (HTTP router/middleware)
- **Database ORM**: GORM with PostgreSQL driver
- **Authentication**: JWT tokens with golang-jwt/jwt
- **OAuth2**: Discord OAuth2 integration
- **Cache**: Redis with go-redis client
- **Testing**: Standard Go testing with testify

## Controller Stack
- **Language**: Go 1.21
- **HTTP Router**: Gorilla Mux
- **Purpose**: Kubernetes controller for GameServer CRDs

## Infrastructure
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Container Registry**: Docker Hub compatible
- **Development Tools**: pgAdmin, Redis Commander

## Common Commands

### Development Setup
```bash
# Quick setup with environment file
make setup

# Start all services
make up

# View logs
make logs

# Check service health
make health
```

### Backend Development
```bash
# Run backend tests with coverage
make test-backend
# OR directly:
cd apps/backend && go test ./... -v

# Run integration tests
cd apps/backend && go test -v ./internal/integration

# Generate coverage report
cd apps/backend && ./run_tests.sh
```

### Frontend Development
```bash
# Run frontend tests
make test-frontend
# OR directly:
cd apps/frontend && npm test

# Development server
cd apps/frontend && npm run dev

# Build for production
cd apps/frontend && npm run build
```

### Database Operations
```bash
# Connect to PostgreSQL
make db-connect

# Backup database
make db-backup

# Restore database
make db-restore
```

### Service Management
```bash
# Individual service logs
make backend-logs
make frontend-logs
make controller-logs

# Rebuild specific services
make backend-build
make frontend-build
make controller-build
```

## Testing Strategy
- **Backend**: Unit tests for handlers, services, middleware + integration tests
- **Frontend**: Component tests with Vitest and Vue Test Utils
- **Coverage**: HTML coverage reports generated for both stacks
- **Environment**: Test-specific configurations and mock services