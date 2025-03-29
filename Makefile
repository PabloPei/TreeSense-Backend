# Makefile for Go project

# Go variables
GO=go
GOFMT=gofmt
LINTER=golangci-lint

# Directories
SRC_DIR=./cmd
BUILD_DIR=./bin

# Default target
all: build

# Build the Go project
build: clean
	$(GO) build -o $(BUILD_DIR)/main $(SRC_DIR)

# Run the Go project
run: build
	$(BUILD_DIR)/main

# Run tests
test:
	$(GO) test ./...

# Lint the code
lint:
	$(LINTER) run

# Format the Go code
fmt:
	$(GOFMT) -w .

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)/*

# Docker build and run
docker-build:
	docker build -t smartspend-backend .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env smartspend-backend

# Create the docker container (build and run)
docker-up: docker-build docker-run

# Install dependencies (go mod tidy)
install:
	$(GO) mod tidy

# Rebuild and run the tests
rebuild-and-test: clean build test

# Help (print this help message)
help:
	@echo "Available commands:"
	@echo "  make all             - Build the project"
	@echo "  make build           - Build the project"
	@echo "  make run             - Run the project"
	@echo "  make test            - Run the tests"
	@echo "  make lint            - Lint the code"
	@echo "  make fmt             - Format the code"
	@echo "  make clean           - Clean the build directory"
	@echo "  make docker-build    - Build the Docker image"
	@echo "  make docker-run      - Run the Docker container"
	@echo "  make docker-up       - Build and run the Docker container"
	@echo "  make install         - Install dependencies"
	@echo "  make rebuild-and-test - Clean, build, and run tests"
	@echo "  make help            - Show this help message"
