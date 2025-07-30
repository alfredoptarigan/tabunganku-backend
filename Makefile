.PHONY: generate-models generate-all-models

# Generate specific tables
generate-models:
	@echo "Generating specific models..."
	@gentool -dsn "host=localhost user=alfredopatriciustarigan password=test dbname=tabanganku port=5432 sslmode=disable TimeZone=UTC" \
		-db postgres \
		-tables "users" \
		-outPath "./pkg/models" \
		-modelPkgName "models" \
		-onlyModel

# Generate all tables
generate-all-models:
	@echo "🚀 Generating models for all tables..."
	@gentool -c ./tools/gen.config.yaml
	@echo "✅ All models generated!"
	@ls -la ./pkg/models/*.go

# Clean and regenerate all
regen-models:
	@echo "🧹 Cleaning existing models..."
	@rm -f ./pkg/models/*.go
	@echo "🚀 Regenerating all models..."
	@make generate-all-models

.PHONY: all build clean run test lint wire migrate seed help

# Default target
all: wire build

# Build the application
build:
	@echo "Building application..."
	go build -o bin/gic-crm cmd/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Generate wire_gen.go files
wire:
	@echo "Generating wire_gen.go files..."
	wire ./pkg/injectors

# Run database migrations
migrate:
	@echo "Running migrations..."
	go run cmd/main.go migrate

# Seed database with initial data
seed:
	@echo "Seeding database..."
	go run cmd/main.go seed

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Default target, run wire and build"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  wire         - Generate wire_gen.go files for dependency injection"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Seed database with initial data"
	@echo "  help         - Show this help message"


# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	GO_ENV=test go test -v -short ./pkg/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v -run Integration ./tests/...

# Run all tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	GO_ENV=test go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in CI/CD environment
test-ci:
	@echo "Running tests for CI/CD..."
	GO_ENV=test go test -v -race -coverprofile=coverage.out ./...

# Setup test environment
setup-test:
	@echo "Setting up test environment..."
	cp example.config.yaml config.yaml
	@echo "Test environment setup complete"