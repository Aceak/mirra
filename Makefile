# Go parameters
BINARY_NAME = mirra
BASE_VERSION := $(shell grep 'const BaseVersion' internal/version/version.go | cut -d'"' -f2)
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_FLAGS = -ldflags="-s -w -X github.com/Aceak/mirra/internal/version.CommitHash=$(COMMIT_HASH)"
GO_CMD = go
GO_BUILD = $(GO_CMD) build $(GO_FLAGS)
GO_TEST = $(GO_CMD) test
GO_FMT = $(GO_CMD) fmt

# Platforms to cross-compile for
# OS: linux, mac (darwin), windows
# Arch: 386 (x86), amd64 (x86_64), arm (ARM 32-bit), arm64 (ARM 64-bit), riscv64 (RISC-V)
PLATFORMS = linux/386 linux/amd64 linux/arm linux/arm64 linux/riscv64 darwin/amd64 darwin/arm64 windows/386 windows/amd64 windows/arm windows/arm64

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	@mkdir -p dist
	@echo "Building $(BINARY_NAME)..."
	@$(GO_BUILD) -o dist/$(BINARY_NAME) ./cmd/server
	@echo "Binary built: dist/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	$(GO_CMD) run main.go

# Run tests
.PHONY: test
test:
	$(GO_TEST) ./...

# Run tests with coverage
.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	@$(GO_TEST) ./... -coverprofile=coverage.out
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO_FMT) ./...
	@echo "Code formatted"

# Lint code
.PHONY: lint
lint:
	@gofmt -d .
	# Add golangci-lint or other linters as needed

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-*.tar.gz $(BINARY_NAME)-*.zip mirra_*.tar.gz mirra_*.zip mirra_* coverage.out coverage.html
	@rm -rf dist/
	@echo "Clean complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO_CMD) mod download
	@$(GO_CMD) mod tidy
	@echo "Dependencies ready"

# Cross-compile for all platforms
.PHONY: cross-build
cross-build: deps
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		if [ "$${GOOS}" = "darwin" ]; then \
			OS_NAME="mac"; \
		else \
			OS_NAME="$${GOOS}"; \
		fi; \
		output_name=dist/mirra_$${OS_NAME}_$${GOARCH}; \
		if [ "$${GOOS}" = "windows" ]; then \
			output_name=$${output_name}.exe; \
		fi; \
		echo "Building $${output_name}..."; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} $(GO_BUILD) -o $${output_name} ./cmd/server; \
	done
	@echo "All builds completed"

# Create distribution packages
.PHONY: dist
dist: clean cross-build
	@echo "Creating distribution packages..."
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		if [ "$${GOOS}" = "darwin" ]; then \
			OS_NAME="mac"; \
		else \
			OS_NAME="$${GOOS}"; \
		fi; \
		output_name=dist/mirra_$${OS_NAME}_$${GOARCH}; \
		archive_name=dist/mirra_$${OS_NAME}_$${GOARCH}; \
		if [ "$${GOOS}" = "windows" ]; then \
			output_name=$${output_name}.exe; \
			cp $${output_name} $(BINARY_NAME).exe; \
			zip -q $${archive_name}.zip $(BINARY_NAME).exe config.example.json README.md; \
			rm -f $(BINARY_NAME).exe; \
		else \
			cp $${output_name} $(BINARY_NAME); \
			tar -czf $${archive_name}.tar.gz $(BINARY_NAME) config.example.json README.md; \
			rm -f $(BINARY_NAME); \
		fi; \
		rm -f $${output_name}; \
		echo "Created $${archive_name}.tar.gz/zip"; \
	done
	@echo "Distribution packages created in dist/"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Default target, builds the project"
	@echo "  build       - Build for current platform"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  cross-build - Cross-compile for multiple platforms"
	@echo "  dist        - Create distribution packages"
	@echo "  help        - Show this help message"