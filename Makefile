# ANSI Editor Framework Makefile

.PHONY: all build test clean editor example examples docs

# Default target
all: build

# Build all components
build:
	@echo "Building ANSI Editor Framework..."
	go build -v ./...

# Build the interactive editor
editor:
	@echo "Building editor..."
	go build -o bin/ansi-editor cmd/editor/main.go

# Build example applications
example:
	@echo "Running example application..."
	go run cmd/example/main.go

# Run example games
examples:
	@echo "Running simple game example..."
	go run examples/simple_game.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf *-export/
	rm -rf example-export/
	rm -rf simple-game-export/
	go clean ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Create binary directory
bin:
	mkdir -p bin

# Build all binaries
binaries: bin
	@echo "Building binaries..."
	go build -o bin/ansi-editor cmd/editor/main.go
	go build -o bin/ansi-example cmd/example/main.go

# Help
help:
	@echo "Available targets:"
	@echo "  all       - Build all components (default)"
	@echo "  build     - Build all packages"
	@echo "  editor    - Build the interactive editor"
	@echo "  example   - Run the example application"
	@echo "  examples  - Run example games"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code (requires golangci-lint)"
	@echo "  binaries  - Build all binaries to bin/"
	@echo "  help      - Show this help"