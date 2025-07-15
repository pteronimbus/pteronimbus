#!/bin/bash

# Test runner script for Pteronimbus backend authentication tests

set -e

echo "🧪 Running Pteronimbus Backend Authentication Tests"
echo "=================================================="

# Change to backend directory
cd "$(dirname "$0")"

# Download dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Run unit tests with coverage
echo ""
echo "🔬 Running unit tests..."
go test -v -race -coverprofile=coverage.out ./internal/handlers ./internal/services ./internal/middleware

# Run integration tests
echo ""
echo "🔗 Running integration tests..."
go test -v -race ./internal/integration

# Generate coverage report
echo ""
echo "📊 Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
echo ""
echo "📈 Coverage Summary:"
go tool cover -func=coverage.out | grep total

echo ""
echo "✅ All tests completed!"
echo "📄 Coverage report generated: coverage.html"
echo ""
echo "Test Results Summary:"
echo "- Unit tests: handlers, services, middleware"
echo "- Integration tests: complete auth flow"
echo "- Coverage report: available in coverage.html"