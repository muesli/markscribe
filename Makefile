# Makefile für GitLab-to-Todoist Exporter

# Variables
BINARY_NAME=markscribe
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin
BINARY_ARM64=$(BINARY_NAME)_arm64

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build configuration
MAIN_FILE=cmd/markscribe/main.go
BUILD_FLAGS=-ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"
BUILD_DIR=./bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Default target
.PHONY: all
all: clean deps lint test build

# Build the binary for current platform
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin build-arm64

.PHONY: build-linux
build-linux:
	@echo "🔨 Building for Linux amd64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_FILE)

.PHONY: build-windows
build-windows:
	@echo "🔨 Building for Windows amd64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) $(MAIN_FILE)

.PHONY: build-darwin
build-darwin:
	@echo "🔨 Building for macOS amd64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_DARWIN) $(MAIN_FILE)

.PHONY: build-arm64
build-arm64:
	@echo "🔨 Building for macOS arm64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_ARM64) $(MAIN_FILE)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f *.md *.csv 2>/dev/null || true
	@echo "✅ Clean completed"

# Test with coverage
.PHONY: test
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "📊 Test coverage:"
	@$(GOCMD) tool cover -func=coverage.out | tail -n 1

# Test coverage report
.PHONY: test-coverage
test-coverage: test
	@echo "🔍 Opening coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

# Lint code
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run; \
	else \
		echo "⚠️  golangci-lint not found, skipping..."; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Download and tidy dependencies
.PHONY: deps
deps:
	@echo "📦 Managing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) verify

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "⬆️  Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Install binary to GOPATH/bin
.PHONY: install
install: build
	@echo "📥 Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(shell $(GOCMD) env GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Installed to $(shell $(GOCMD) env GOPATH)/bin/$(BINARY_NAME)"

# Development helpers
.PHONY: dev
dev: fmt lint test build

# Quick development cycle
.PHONY: dev-quick
dev-quick: fmt build

# Run the application
.PHONY: run
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with debug mode
.PHONY: debug
debug: build
	@echo "🔍 Running debug mode..."
	@$(BUILD_DIR)/$(BINARY_NAME) --debug

# Run with help
.PHONY: help-app
help-app: build
	@$(BUILD_DIR)/$(BINARY_NAME) --help

# Run examples
.PHONY: run-markdown
run-markdown: build
	@echo "📝 Running Markdown export example..."
	@$(BUILD_DIR)/$(BINARY_NAME) --output-file example_export.md --verbose

.PHONY: run-todoist
run-todoist: build
	@echo "📋 Running Todoist export example..."
	@$(BUILD_DIR)/$(BINARY_NAME) --todoist --verbose

# Setup development environment
.PHONY: setup
setup:
	@echo "🛠️  Setting up development environment..."
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			cp .env.example .env; \
			echo "📋 Created .env from .env.example - please edit with your tokens"; \
		else \
			echo "Creating basic .env file..."; \
			echo "# GitLab Configuration" > .env; \
			echo "GITLAB_TOKEN=glpat-your-token-here" >> .env; \
			echo "PROJECT_PATH=your-username/your-project" >> .env; \
			echo "" >> .env; \
			echo "# Todoist Configuration" >> .env; \
			echo "TODOIST_TOKEN=your-todoist-token-here" >> .env; \
			echo "TODOIST_PROJECT=GitLab Issues" >> .env; \
			echo "TODOIST_API=false" >> .env; \
			echo "" >> .env; \
			echo "# Optional" >> .env; \
			echo "#MILESTONE_TITLE=v1.0.0" >> .env; \
			echo "#VERBOSE=true" >> .env; \
		fi; \
	else \
		echo "⚠️  .env already exists"; \
	fi
	@echo "📦 Installing dependencies..."
	@$(MAKE) deps
	@echo "✅ Setup completed!"

# Create release archives
.PHONY: release
release: clean build-all
	@echo "📦 Creating release archives v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)/releases
	@mkdir -p $(BUILD_DIR)/temp
	
	# Linux
	@cp $(BUILD_DIR)/$(BINARY_UNIX) $(BUILD_DIR)/temp/$(BINARY_NAME)
	@tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BUILD_DIR)/temp $(BINARY_NAME)

	# Darwin amd64
	@cp $(BUILD_DIR)/$(BINARY_DARWIN) $(BUILD_DIR)/temp/$(BINARY_NAME)
	@tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(BUILD_DIR)/temp $(BINARY_NAME)

	# Darwin arm64
	@cp $(BUILD_DIR)/$(BINARY_ARM64) $(BUILD_DIR)/temp/$(BINARY_NAME)
	@tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(BUILD_DIR)/temp $(BINARY_NAME)
	
	# Windows
	@cd $(BUILD_DIR) && zip -j releases/$(BINARY_NAME)-windows-amd64.zip $(BINARY_WINDOWS)
	
	@echo "📦 Release archives created in $(BUILD_DIR)/releases/"

# Validate project structure
.PHONY: validate
validate:
	@echo "🔍 Validating project structure..."
	@test -f $(MAIN_FILE) || (echo "❌ $(MAIN_FILE) not found" && exit 1)
	@test -d internal/ || (echo "❌ internal/ directory not found" && exit 1)
	@test -f go.mod || (echo "❌ go.mod not found" && exit 1)
	@echo "✅ Project structure valid"

# Security audit
.PHONY: audit
audit:
	@echo "🛡️  Running security audit..."
	@$(GOCMD) list -json -m all | nancy sleuth
	@govulncheck ./...

# Show application info
.PHONY: info
info:
	@echo "📊 Build Information:"
	@echo "  Binary Name:  $(BINARY_NAME)"
	@echo "  Version:      $(VERSION)"
	@echo "  Build Dir:    $(BUILD_DIR)"
	@echo "  Main File:    $(MAIN_FILE)"
	@echo "  Go Version:   $(shell $(GOCMD) version)"
	@echo "  OS/Arch:      $(shell $(GOCMD) env GOOS)/$(shell $(GOCMD) env GOARCH)"
	@echo "  Build Time:   $(BUILD_TIME)"

# Show help
.PHONY: help
help:
	@echo "🛠️  GitLab-to-Todoist Exporter - Available Targets:"
	@echo ""
	@echo "📦 BUILD:"
	@echo "  build        - Build binary for current platform"
	@echo "  build-all    - Build binaries for all platforms (Linux, Windows, macOS)"
	@echo "  build-linux  - Build binary for Linux amd64"
	@echo "  build-windows- Build binary for Windows amd64"
	@echo "  build-darwin - Build binary for macOS amd64"
	@echo "  build-arm64  - Build binary for macOS arm64"
	@echo ""
	@echo "🧪 TESTING:"
	@echo "  test         - Run tests with coverage"
	@echo "  test-coverage- Generate HTML coverage report"
	@echo ""
	@echo "🔧 DEVELOPMENT:"
	@echo "  dev          - Full development cycle (fmt, lint, test, build)"
	@echo "  dev-quick    - Quick cycle (fmt, build)"
	@echo "  fmt          - Format source code"
	@echo "  lint         - Run linter"
	@echo "  setup        - Setup development environment (.env, deps)"
	@echo ""
	@echo "📦 DEPENDENCIES:"
	@echo "  deps         - Download and verify dependencies"
	@echo "  deps-update  - Update all dependencies"
	@echo ""
	@echo "🚀 RUN:"
	@echo "  run          - Build and run with .env configuration"
	@echo "  debug        - Run in debug mode"
	@echo "  help-app     - Show application help"
	@echo "  run-markdown - Example: Export to Markdown"
	@echo "  run-todoist  - Example: Export to Todoist"
	@echo ""
	@echo "🚀 DISTRIBUTION:"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  release      - Create release archives for all platforms"
	@echo ""
	@echo "🔍 UTILITIES:"
	@echo "  clean        - Remove build artifacts"
	@echo "  validate     - Validate project structure"
	@echo "  audit        - Run security audit"
	@echo "  info         - Show build information"
	@echo "  help         - Show this help"
