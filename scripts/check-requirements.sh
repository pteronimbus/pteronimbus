#!/bin/bash

# Check Requirements Script for Pteronimbus Docker Setup

set -e

echo "🔍 Checking system requirements for Pteronimbus..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | cut -d',' -f1)
    echo -e "${GREEN}✅ Docker is installed${NC} (version: $DOCKER_VERSION)"
else
    echo -e "${RED}❌ Docker is not installed${NC}"
    echo "Please install Docker from: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version | cut -d' ' -f3 | cut -d',' -f1)
    echo -e "${GREEN}✅ Docker Compose is installed${NC} (version: $COMPOSE_VERSION)"
elif docker compose version &> /dev/null; then
    COMPOSE_VERSION=$(docker compose version --short)
    echo -e "${GREEN}✅ Docker Compose (plugin) is installed${NC} (version: $COMPOSE_VERSION)"
else
    echo -e "${RED}❌ Docker Compose is not installed${NC}"
    echo "Please install Docker Compose from: https://docs.docker.com/compose/install/"
    exit 1
fi

# Check if Docker daemon is running
if docker info &> /dev/null; then
    echo -e "${GREEN}✅ Docker daemon is running${NC}"
else
    echo -e "${RED}❌ Docker daemon is not running${NC}"
    echo "Please start Docker daemon"
    exit 1
fi

# Check available disk space (warn if less than 5GB)
AVAILABLE_SPACE=$(df -BG . | tail -1 | awk '{print $4}' | sed 's/G//')
if [ "$AVAILABLE_SPACE" -lt 5 ]; then
    echo -e "${YELLOW}⚠️  Warning: Low disk space (${AVAILABLE_SPACE}GB available)${NC}"
    echo "Docker images and containers may require several GB of space"
else
    echo -e "${GREEN}✅ Sufficient disk space available${NC} (${AVAILABLE_SPACE}GB)"
fi

# Check if required ports are available
echo ""
echo "🔍 Checking port availability..."

check_port() {
    local port=$1
    local service=$2
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Port $port is already in use${NC} (needed for $service)"
        echo "   You may need to stop the service using this port or change the port in docker-compose.yml"
    else
        echo -e "${GREEN}✅ Port $port is available${NC} ($service)"
    fi
}

check_port 3000 "Frontend"
check_port 8080 "Backend API"
check_port 8081 "Controller"
check_port 5433 "PostgreSQL"
check_port 6380 "Redis"
check_port 8082 "Redis Commander"
check_port 8083 "pgAdmin"

# Check if .env file exists
echo ""
echo "🔍 Checking configuration..."

if [ -f ".env" ]; then
    echo -e "${GREEN}✅ .env file exists${NC}"
    
    # Check if Discord credentials are set
    if grep -q "DISCORD_CLIENT_ID=your_discord_client_id" .env; then
        echo -e "${YELLOW}⚠️  Discord OAuth2 credentials not configured${NC}"
        echo "   Please edit .env file and add your Discord application credentials"
    else
        echo -e "${GREEN}✅ Discord OAuth2 credentials appear to be configured${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  .env file not found${NC}"
    echo "   Run 'make setup-env' or 'cp .env.example .env' to create it"
fi

echo ""
echo "🎉 System check complete!"
echo ""
echo "Next steps:"
echo "1. If you don't have a .env file: run 'make setup-env'"
echo "2. Configure Discord OAuth2 credentials in .env file"
echo "3. Start the services: run 'make up' or 'docker-compose up -d'"
echo "4. Check service health: run 'make health'"
echo ""
echo "For more information, see DOCKER_SETUP.md"