# Test runner script for Pteronimbus backend authentication tests

Write-Host "🧪 Running Pteronimbus Backend Authentication Tests" -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# Change to backend directory
Set-Location $PSScriptRoot

# Download dependencies
Write-Host ""
Write-Host "📦 Installing dependencies..." -ForegroundColor Yellow
go mod tidy

# Run unit tests with coverage
Write-Host ""
Write-Host "🔬 Running unit tests..." -ForegroundColor Yellow
go test -v -race -coverprofile=coverage.out ./internal/handlers ./internal/services ./internal/middleware

# Run integration tests
Write-Host ""
Write-Host "🔗 Running integration tests..." -ForegroundColor Yellow
go test -v -race ./internal/integration

# Generate coverage report
Write-Host ""
Write-Host "📊 Generating coverage report..." -ForegroundColor Yellow
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
Write-Host ""
Write-Host "📈 Coverage Summary:" -ForegroundColor Green
go tool cover -func=coverage.out | Select-String "total"

Write-Host ""
Write-Host "✅ All tests completed!" -ForegroundColor Green
Write-Host "📄 Coverage report generated: coverage.html" -ForegroundColor Green
Write-Host ""
Write-Host "Test Results Summary:" -ForegroundColor White
Write-Host "- Unit tests: handlers, services, middleware" -ForegroundColor White
Write-Host "- Integration tests: complete auth flow" -ForegroundColor White
Write-Host "- Coverage report: available in coverage.html" -ForegroundColor White