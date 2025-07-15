# Pteronimbus Docker Setup

This Docker Compose setup provides a complete development environment for Pteronimbus with all required services.

## Services Included

### Core Services
- **Frontend** (Nuxt.js) - Port 3000
- **Backend** (Go API) - Port 8080  
- **Controller** (Go Controller) - Port 8081
- **PostgreSQL** (Database) - Port 5433
- **Redis** (Cache/Sessions) - Port 6380

### Management Tools (Optional)
- **pgAdmin** (PostgreSQL Management) - Port 8083
- **Redis Commander** (Redis Management) - Port 8082

## Quick Start

1. **Clone and navigate to the project:**
   ```bash
   git clone <repository-url>
   cd pteronimbus
   ```

2. **Set up Discord OAuth2 (Required):**
   ```bash
   cp .env.example .env
   # Edit .env and add your Discord OAuth2 credentials
   ```

3. **Start all services:**
   ```bash
   docker-compose up -d
   ```

4. **Check service status:**
   ```bash
   docker-compose ps
   ```

## Discord OAuth2 Setup

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to OAuth2 settings
4. Add redirect URL: `http://localhost:8080/auth/callback`
5. Copy Client ID and Client Secret to your `.env` file

## Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Main application |
| Backend API | http://localhost:8080 | REST API |
| Controller | http://localhost:8081 | Kubernetes controller |
| pgAdmin | http://localhost:8083 | Database management |
| Redis Commander | http://localhost:8082 | Redis management |

## Default Credentials

### Database (PostgreSQL)
- **Host:** localhost:5433
- **Database:** pteronimbus
- **Username:** postgres
- **Password:** postgres123

### Redis
- **Host:** localhost:6380
- **Password:** redis123

### pgAdmin
- **Email:** admin@pteronimbus.local
- **Password:** admin123

### Redis Commander
- **Username:** admin
- **Password:** admin123

## Common Commands

### Start all services
```bash
docker-compose up -d
```

### Stop all services
```bash
docker-compose down
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
```

### Rebuild and restart a service
```bash
docker-compose up -d --build backend
```

### Access service shell
```bash
# Backend container
docker-compose exec backend sh

# Database
docker-compose exec postgres psql -U postgres -d pteronimbus
```

### Reset all data
```bash
docker-compose down -v
docker-compose up -d
```

## Development Workflow

### Making Code Changes

1. **Backend changes:**
   ```bash
   docker-compose up -d --build backend
   ```

2. **Frontend changes:**
   ```bash
   docker-compose up -d --build frontend
   ```

3. **Controller changes:**
   ```bash
   docker-compose up -d --build controller
   ```

### Database Operations

1. **Connect to PostgreSQL:**
   ```bash
   docker-compose exec postgres psql -U postgres -d pteronimbus
   ```

2. **Run SQL scripts:**
   ```bash
   docker-compose exec postgres psql -U postgres -d pteronimbus -f /path/to/script.sql
   ```

3. **Backup database:**
   ```bash
   docker-compose exec postgres pg_dump -U postgres pteronimbus > backup.sql
   ```

### Redis Operations

1. **Connect to Redis:**
   ```bash
   docker-compose exec redis redis-cli -a redis123
   ```

2. **Monitor Redis:**
   ```bash
   docker-compose exec redis redis-cli -a redis123 monitor
   ```

## Troubleshooting

### Service won't start
```bash
# Check logs
docker-compose logs service-name

# Check if ports are in use
netstat -tulpn | grep :PORT_NUMBER
```

### Database connection issues
```bash
# Check if PostgreSQL is ready
docker-compose exec postgres pg_isready -U postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### Redis connection issues
```bash
# Test Redis connection
docker-compose exec redis redis-cli -a redis123 ping
```

### Build issues
```bash
# Clean build
docker-compose down
docker system prune -f
docker-compose up -d --build
```

## Environment Variables

You can override default settings by creating a `.env` file:

```bash
# Discord OAuth2
DISCORD_CLIENT_ID=your_client_id
DISCORD_CLIENT_SECRET=your_client_secret

# Database
POSTGRES_PASSWORD=custom_password

# Redis  
REDIS_PASSWORD=custom_password
```

## Production Considerations

This setup is designed for development. For production:

1. Use proper secrets management
2. Enable SSL/TLS
3. Use production-grade database settings
4. Implement proper backup strategies
5. Use container orchestration (Kubernetes)
6. Set up monitoring and logging

## Health Checks

All services include health checks. You can monitor them with:

```bash
docker-compose ps
```

Services will show as "healthy" when ready to accept connections.

## Volumes

Persistent data is stored in Docker volumes:
- `postgres_data` - Database files
- `redis_data` - Redis persistence
- `pgadmin_data` - pgAdmin settings

To backup volumes:
```bash
docker run --rm -v pteronimbus_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .
```