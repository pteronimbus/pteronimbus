# Project Structure

## Repository Organization
Pteronimbus follows a monorepo structure with clear separation of concerns across multiple applications and shared resources.

## Root Level Structure
```
├── apps/                    # Application services
├── charts/                  # Helm deployment charts
├── config/                  # Kubernetes configurations (CRDs, RBAC, samples)
├── deployments/             # Deployment examples and configurations
├── docs/                    # Docusaurus documentation site
├── internal/                # Shared internal packages
├── proto/                   # Protocol buffer definitions
├── scripts/                 # Development and setup scripts
├── docker-compose.yml       # Local development orchestration
├── Makefile                 # Development commands and automation
└── .kiro/                   # Kiro AI assistant configuration
```

## Applications Structure (`apps/`)

### Backend (`apps/backend/`)
```
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/          # HTTP middleware (auth, CORS, tenant)
│   ├── models/              # Data models and structures
│   └── services/            # Business logic services
├── coverage.out             # Test coverage output
├── go.mod                   # Go module definition
└── run_tests.sh             # Test execution script
```

### Frontend (`apps/frontend/`)
```
├── components/              # Vue components
├── composables/             # Vue composables (useAuth, useTenant)
├── middleware/              # Nuxt middleware (auth, guest, tenant)
├── pages/                   # File-based routing pages
├── plugins/                 # Nuxt plugins
├── server/api/              # Server-side API routes
├── test/                    # Test files mirroring src structure
├── assets/css/              # Stylesheets
├── i18n/                    # Internationalization
├── public/                  # Static assets
├── nuxt.config.ts           # Nuxt configuration
└── package.json             # Node.js dependencies
```

### Controller (`apps/controller/`)
```
├── cmd/server/              # Application entry point
├── internal/handlers/       # HTTP handlers
├── api/                     # API definitions (placeholder)
├── config/                  # Configuration (placeholder)
├── controllers/             # Kubernetes controllers (placeholder)
└── go.mod                   # Go module definition
```

## Configuration Structure

### Kubernetes Config (`config/`)
```
├── crd/                     # Custom Resource Definitions
├── rbac/                    # Role-Based Access Control
└── samples/                 # Example configurations
```

### Helm Charts (`charts/`)
```
├── pteronimbus/             # Main application chart
└── controller/              # Controller-specific chart
```

## Documentation (`docs/`)
- **Framework**: Docusaurus for documentation site
- **Structure**: Standard Docusaurus layout with docs/, blog/, src/
- **Content**: Architecture, API, user guides, tutorials

## Development Files

### Environment Configuration
- `.env.example` - Template for environment variables
- `.env` - Local environment variables (gitignored)
- Individual app `.env.example` files

### Docker Configuration
- `docker-compose.yml` - Multi-service development environment
- Individual `Dockerfile` per application
- Health checks and service dependencies

### Testing Structure
- **Backend**: Tests alongside source in `internal/` subdirectories
- **Frontend**: Dedicated `test/` directory mirroring component structure
- **Integration**: Separate integration test directories
- **Coverage**: HTML reports generated in `coverage/` directories

## Naming Conventions

### Files and Directories
- **Go**: Snake_case for files, PascalCase for types
- **TypeScript/Vue**: kebab-case for files, PascalCase for components
- **Tests**: `*.test.ts` (frontend), `*_test.go` (backend)
- **Config**: kebab-case for configuration files

### Code Organization
- **Handlers**: HTTP endpoint handlers grouped by domain
- **Services**: Business logic separated from HTTP concerns
- **Models**: Data structures and database entities
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)
- **Composables**: Reusable Vue composition functions

## Import Patterns
- **Backend**: Internal packages use full module path
- **Frontend**: Nuxt auto-imports for composables, components, utils
- **Shared**: Common types and utilities in `internal/` packages