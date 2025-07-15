# Pteronimbus Backend

This is the Go backend for Pteronimbus that handles Discord OAuth2 authentication, JWT token management, and API endpoints.

## Features

- Discord OAuth2 authentication flow
- JWT token generation and validation
- Refresh token mechanism with Redis storage
- Authentication middleware for API routes
- CORS support for frontend integration
- Health check endpoints

## Setup

### Prerequisites

- Go 1.21+
- Redis server
- Discord OAuth2 application

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Server Configuration
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000

# Discord OAuth2 Configuration
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DISCORD_REDIRECT_URL=http://localhost:8080/auth/callback

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ISSUER=pteronimbus

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Discord OAuth2 Setup

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to OAuth2 settings
4. Add redirect URL: `http://localhost:8080/auth/callback`
5. Copy Client ID and Client Secret to your environment variables

### Running the Server

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

## API Endpoints

### Authentication

- `GET /auth/login` - Get Discord OAuth2 authorization URL
- `GET /auth/callback` - Handle Discord OAuth2 callback
- `POST /auth/refresh` - Refresh access token
- `GET /auth/me` - Get current user info (requires auth)
- `POST /auth/logout` - Logout and invalidate session (requires auth)

### Health Checks

- `GET /health` - General health check
- `GET /healthz` - Kubernetes-style health check
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe

### Protected API

- `GET /api/test` - Test protected endpoint (requires auth)

## Authentication Flow

1. Frontend calls `/auth/login` to get Discord OAuth2 URL
2. User is redirected to Discord for authorization
3. Discord redirects back to `/auth/callback` with authorization code
4. Backend exchanges code for Discord access token
5. Backend creates user session and generates JWT tokens
6. Frontend receives access token and refresh token
7. Frontend uses access token for API requests
8. When access token expires, frontend uses refresh token to get new access token

## Architecture

- **Config**: Environment-based configuration management
- **Services**: Business logic (Auth, Discord, JWT, Redis)
- **Handlers**: HTTP request handlers
- **Middleware**: Authentication and CORS middleware
- **Models**: Data structures and types

## Security Features

- CSRF protection with state parameter in OAuth2 flow
- JWT tokens with expiration
- Secure session storage in Redis
- CORS configuration
- Bearer token authentication
- Automatic token refresh