.PHONY: all test test-cover build run clean fmt vet

# Variables
BINARY_NAME_INGEST=grextor-ingest
BINARY_NAME_QUERY=grextor-query
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

all: fmt vet test build

# Formatting
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet
vet:
	@echo "Vetting code..."
	@go vet ./...

# Testing
test:
	@echo "Running tests..."
	@go test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

# Building
build:
	@echo "Building binaries..."
	@go build -o $(BINARY_NAME_INGEST) ./cmd/ingest
	@go build -o $(BINARY_NAME_QUERY) ./cmd/query

# Running (Example: run query by default, or provide target)
run: build
	@echo "Run 'minikube service...' or specific binary directly"
	@echo "Example: ./$(BINARY_NAME_QUERY)"

# Cleaning
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_NAME_INGEST)
	@rm -f $(BINARY_NAME_QUERY)
	@rm -f coverage.out
