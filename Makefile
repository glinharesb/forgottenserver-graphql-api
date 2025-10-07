.PHONY: help test test-verbose test-coverage run build clean docker-build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out | tail -1
	@echo "Coverage report saved to coverage.out"
	@echo "To view HTML report: go tool cover -html=coverage.out"

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test ./... -v -cover

test-models: ## Run only model tests
	@echo "Running model tests..."
	@go test ./internal/models -v -cover

test-graph: ## Run only graph/resolver tests
	@echo "Running graph tests..."
	@go test ./internal/graph -v -cover

run: ## Run the server
	@echo "Starting GraphQL server..."
	@go run cmd/server/main.go

build: ## Build the binary
	@echo "Building..."
	@go build -o bin/api-graphql cmd/server/main.go
	@echo "Binary created at bin/api-graphql"

clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Done!"

generate: ## Run gqlgen code generation
	@echo "Generating GraphQL code..."
	@go run github.com/99designs/gqlgen generate
	@echo "Done!"

tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	@go mod tidy
	@echo "Done!"

# Docker commands

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t forgottenserver-graphql-api:latest .
	@echo "Done!"
