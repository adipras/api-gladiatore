.PHONY: help build run clean frontend-dev frontend-build backend-dev docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the backend server"
	@echo "  make run            - Run the complete application"
	@echo "  make frontend-dev   - Run frontend in development mode"
	@echo "  make frontend-build - Build frontend for production"
	@echo "  make backend-dev    - Run backend in development mode"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run in Docker container"
	@echo "  make test           - Run tests"

# Build backend
build:
	go build -o bin/api-gladiatore cmd/server/main.go
	@echo "✅ Backend built successfully"

# Run complete application
run: frontend-build
	go run cmd/server/main.go

# Frontend development
frontend-dev:
	cd frontend && npm start

# Frontend production build
frontend-build:
	cd frontend && npm run build
	@echo "✅ Frontend built successfully"

# Backend development (with hot reload using air if installed)
backend-dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		go run cmd/server/main.go; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf generated/*
	rm -rf frontend/build/
	@echo "✅ Cleaned build artifacts"

# Docker build
docker-build:
	docker build -t api-gladiatore .

# Docker run
docker-run:
	docker run -p 8080:8080 api-gladiatore

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod download
	cd frontend && npm install
	@echo "✅ Dependencies installed"

# Database setup
db-setup:
	@echo "Creating database if not exists..."
	mysql -u root -e "CREATE DATABASE IF NOT EXISTS api_gladiatore;"
	@echo "✅ Database setup complete"

# Full setup (install deps, setup db, build)
setup: deps db-setup frontend-build build
	@echo "✅ Setup complete! Run 'make run' to start the application"