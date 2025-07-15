# Pteronimbus Development Makefile

.PHONY: help up down logs build clean test setup-env

# Default target
help: ## Show this help message
	@echo "Pteronimbus Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose Commands
up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## Show logs for all services
	docker-compose logs -f

build: ## Build and start all services
	docker-compose up -d --build

clean: ## Stop services and remove volumes (WARNING: This will delete all data!)
	docker-compose down -v
	docker system prune -f

restart: ## Restart all services
	docker-compose restart

# Individual Service Commands
backend-logs: ## Show backend logs
	docker-compose logs -f backend

frontend-logs: ## Show frontend logs
	docker-compose logs -f frontend

controller-logs: ## Show controller logs
	docker-compose logs -f controller

backend-build: ## Rebuild and restart backend
	docker-compose up -d --build backend

frontend-build: ## Rebuild and restart frontend
	docker-compose up -d --build frontend

controller-build: ## Rebuild and restart controller
	docker-compose up -d --build controller

# Database Commands
db-connect: ## Connect to PostgreSQL database
	docker-compose exec postgres psql -U postgres -d pteronimbus

db-backup: ## Backup database to backup.sql
	docker-compose exec postgres pg_dump -U postgres pteronimbus > backup.sql

db-restore: ## Restore database from backup.sql
	docker-compose exec -T postgres psql -U postgres -d pteronimbus < backup.sql

# Redis Commands
redis-connect: ## Connect to Redis CLI
	docker-compose exec redis redis-cli -a redis123

redis-monitor: ## Monitor Redis commands
	docker-compose exec redis redis-cli -a redis123 monitor

# Development Commands
test-backend: ## Run backend tests
	cd apps/backend && go test ./... -v

test-frontend: ## Run frontend tests
	cd apps/frontend && npm test

setup-env: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please edit it with your Discord OAuth2 credentials."; \
	else \
		echo ".env file already exists."; \
	fi

# Status Commands
status: ## Show status of all services
	docker-compose ps

health: ## Check health of all services
	@echo "Checking service health..."
	@curl -s http://localhost:3000/ > /dev/null && echo "✅ Frontend: OK" || echo "❌ Frontend: DOWN"
	@curl -s http://localhost:8080/health > /dev/null && echo "✅ Backend: OK" || echo "❌ Backend: DOWN"
	@curl -s http://localhost:8081/health > /dev/null && echo "✅ Controller: OK" || echo "❌ Controller: DOWN"

# Quick Setup
setup: setup-env up ## Quick setup: create .env and start all services
	@echo ""
	@echo "🚀 Pteronimbus is starting up!"
	@echo ""
	@echo "Services will be available at:"
	@echo "  Frontend:        http://localhost:3000"
	@echo "  Backend API:     http://localhost:8080"
	@echo "  Controller:      http://localhost:8081"
	@echo "  pgAdmin:         http://localhost:8083"
	@echo "  Redis Commander: http://localhost:8082"
	@echo ""
	@echo "⚠️  Don't forget to configure your Discord OAuth2 credentials in .env"
	@echo ""
	@echo "Run 'make logs' to see startup logs"
	@echo "Run 'make health' to check service status"

# Development helpers
dev-backend: ## Start only backend development stack (backend + db + redis)
	docker-compose up -d postgres redis backend

dev-frontend: ## Start only frontend development stack (frontend + backend + db + redis)
	docker-compose up -d postgres redis backend frontend

# Cleanup commands
clean-images: ## Remove unused Docker images
	docker image prune -f

clean-all: ## Full cleanup (containers, volumes, images, networks)
	docker-compose down -v
	docker system prune -af
	docker volume prune -f