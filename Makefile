.PHONY: help build run test clean docker-build docker-up docker-down deps

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build all services"
	@echo "  build-producer  - Build producer service"
	@echo "  build-consumer  - Build consumer service"
	@echo "  build-notification - Build notification service"
	@echo "  run-producer    - Run producer service locally"
	@echo "  run-consumer    - Run consumer service locally"
	@echo "  run-notification - Run notification service locally"
	@echo "  test            - Run all tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Download dependencies"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-up       - Start all services with Docker Compose"
	@echo "  docker-down     - Stop all services"
	@echo "  docker-logs     - Show Docker logs"
	@echo "  swagger         - Generate Swagger documentation"

# Build targets
build: build-producer build-consumer build-notification

build-producer:
	@echo "Building producer service..."
	go build -o bin/producer-service ./cmd/producer-service

build-consumer:
	@echo "Building consumer service..."
	go build -o bin/consumer-service ./cmd/consumer-service

build-notification:
	@echo "Building notification service..."
	go build -o bin/notification-service ./cmd/notification-service

# Run targets (for local development)
run-producer:
	@echo "Starting producer service on port 3000..."
	PORT=3000 go run ./cmd/producer-service

run-consumer:
	@echo "Starting consumer service on port 3001..."
	PORT=3001 go run ./cmd/consumer-service

run-notification:
	@echo "Starting notification service on port 3002..."
	PORT=3002 go run ./cmd/notification-service

# Test targets
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Utility targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Docker targets
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

docker-restart: docker-down docker-up

# Development targets
dev-kafka:
	@echo "Starting only Kafka and Zookeeper..."
	docker-compose up -d zookeeper kafka kafka-setup

dev-stop-kafka:
	@echo "Stopping Kafka and Zookeeper..."
	docker-compose stop zookeeper kafka kafka-setup

# Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/producer-service/main.go -o cmd/producer-service/docs
	swag init -g cmd/consumer-service/main.go -o cmd/consumer-service/docs
	swag init -g cmd/notification-service/main.go -o cmd/notification-service/docs

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create directories
init:
	@echo "Creating directories..."
	mkdir -p bin/
