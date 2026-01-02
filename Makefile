.PHONY: help run build test clean migrate-up migrate-down migrate-status swagger docker-up docker-down install lint

# Default target
help:
	@echo "VNalo Backend - Available Commands:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback last migration"
	@echo "  make migrate-status - Show migration status"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make install       - Install dependencies"
	@echo "  make lint          - Run linter"

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	@echo "Building VNalo Backend..."
	go build -o bin/vnalo-api cmd/api/main.go
	go build -o bin/vnalo-migrate cmd/migrate/main.go
	@echo "Build complete! Binaries in ./bin/"

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf cmd/swag/docs/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# PostgreSQL Migrations (using Goose + GORM)
migrate-up:
	go run cmd/migrate/main.go -cmd=up

migrate-down:
	go run cmd/migrate/main.go -cmd=down

migrate-status:
	go run cmd/migrate/main.go -cmd=status

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o cmd/swag/docs
	@echo "Swagger docs generated at ./cmd/swag/docs/"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development mode (with hot reload)
dev:
	@echo "Starting development mode..."
	air

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Mod tidy
tidy:
	go mod tidy

# PostgreSQL specific migrations
migrate-postgres-up:
	goose -dir migrations/postgres postgres "host=localhost port=5432 user=postgres password=postgres dbname=vnalo_db sslmode=disable" up

migrate-postgres-down:
	goose -dir migrations/postgres postgres "host=localhost port=5432 user=postgres password=postgres dbname=vnalo_db sslmode=disable" down

migrate-postgres-status:
	goose -dir migrations/postgres postgres "host=localhost port=5432 user=postgres password=postgres dbname=vnalo_db sslmode=disable" status

# Cassandra schema management
cassandra-init:
	docker cp migrations/cassandra/schema.cql vnalo_cassandra:/schema.cql
	docker exec -it vnalo_cassandra cqlsh -f /schema.cql

cassandra-reset:
	docker exec -it vnalo_cassandra cqlsh -e "DROP KEYSPACE IF EXISTS vnalo_chat;"
	$(MAKE) cassandra-init

cassandra-shell:
	docker exec -it vnalo_cassandra cqlsh

# Complete setup for development
setup-dev:
	@echo "Setting up development environment..."
	$(MAKE) docker-up
	@echo "Waiting for services to be ready..."
	@sleep 60
	$(MAKE) cassandra-init
	$(MAKE) migrate-postgres-up
	@echo "Development environment ready!"

# Database shells
postgres-shell:
	docker exec -it vnalo_postgres psql -U postgres -d vnalo_db

redis-shell:
	docker exec -it vnalo_redis redis-cli

