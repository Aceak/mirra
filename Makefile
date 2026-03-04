# Go parameters
BINARY_NAME = fileserver
BASE_VERSION := $(shell grep 'const baseVersion' version.go | cut -d'"' -f2)
BUILD_ID ?= $(shell date -u '+%s')
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION = $(BASE_VERSION).$(BUILD_ID)
GO_FLAGS = -ldflags="-s -w -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH) -X main.buildTime=$(BUILD_TIME)"
GO_CMD = go
GO_BUILD = $(GO_CMD) build $(GO_FLAGS)
GO_TEST = $(GO_CMD) test
GO_FMT = $(GO_CMD) fmt

# Platforms to cross-compile for
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 freebsd/amd64

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	@mkdir -p dist
	$(GO_BUILD) -o dist/$(BINARY_NAME) .
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
	$(GO_TEST) ./... -coverprofile=coverage.out
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	$(GO_FMT) ./...

# Lint code
.PHONY: lint
lint:
	gofmt -d .
	# Add golangci-lint or other linters as needed

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*.tar.gz $(BINARY_NAME)-*.zip coverage.out coverage.html
	rm -rf dist/

# Install dependencies
.PHONY: deps
deps:
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

# Cross-compile for all platforms
.PHONY: cross-build
cross-build: deps
	@echo "Building for multiple platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=dist/$(BINARY_NAME)-$${GOOS}-$${GOARCH}; \
		if [ "$${GOOS}" = "windows" ]; then \
			output_name=$${output_name}.exe; \
		fi; \
		echo "Building $${output_name}..."; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} $(GO_BUILD) -o $${output_name} .; \
	done

# Create distribution packages
.PHONY: dist
dist: clean cross-build
	@echo "Creating distribution packages..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=dist/$(BINARY_NAME)-$${GOOS}-$${GOARCH}; \
		archive_name=dist/$(BINARY_NAME)-$(VERSION)-$${GOOS}-$${GOARCH}; \
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
	@echo "Distribution packages created in dist/ directory"

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

# Add version information to the binary
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_HASH)"
	@echo "Build Time: $(BUILD_TIME)"