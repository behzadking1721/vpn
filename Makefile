# Makefile for VPN Client

# Variables
BINARY_NAME = vpn-client
BINARY_DIR = bin
SRC_DIR = src
UI_DIR = ui
DOCS_DIR = docs
DEMO_DIR = $(SRC_DIR)/demo

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Default target
all: deps build

# Build the application
build:
	@echo "Building VPN Client..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "Build completed successfully!"

# Build demo application
build-demo:
	@echo "Building demo application..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/vpn-demo $(DEMO_DIR)/main.go
	@echo "Demo build completed successfully!"

# Build for different platforms
build-all: build-windows build-linux build-macos

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BINARY_DIR)/windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/windows/$(BINARY_NAME).exe $(SRC_DIR)/main.go
	@echo "Windows build completed!"

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BINARY_DIR)/linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/linux/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "Linux build completed!"

build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BINARY_DIR)/macos
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/macos/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "macOS build completed!"

# Run the application in normal mode
run:
	@echo "Running VPN Client..."
	cd $(SRC_DIR) && $(GOCMD) run main.go

# Run the application in API mode
run-api:
	@echo "Running VPN Client API Server..."
	cd $(SRC_DIR) && $(GOCMD) run main.go --api

# Run the application in CLI mode
run-cli:
	@echo "Running VPN Client CLI..."
	cd $(SRC_DIR) && $(GOCMD) run main.go --cli

# Run demo application
run-demo:
	@echo "Running demo application..."
	cd $(DEMO_DIR) && $(GOCMD) run main.go

# Clean build files
clean:
	@echo "Cleaning build files..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	@echo "Clean completed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	cd $(SRC_DIR) && $(GOMOD) tidy
	@echo "Dependencies installed!"

# Run tests
test:
	@echo "Running tests..."
	cd $(SRC_DIR) && $(GOTEST) -v ./...

# Run protocol tests
test-protocols:
	@echo "Running protocol tests..."
	cd $(SRC_DIR) && $(GOTEST) -v ./protocols

# Run manager tests
test-managers:
	@echo "Running manager tests..."
	cd $(SRC_DIR) && $(GOTEST) -v ./managers

# Generate documentation
docs:
	@echo "Generating documentation..."
	@mkdir -p $(DOCS_DIR)/generated
	@echo "Documentation generated in $(DOCS_DIR)/generated"

# Help
help:
	@echo "VPN Client Makefile"
	@echo "=================="
	@echo "Available targets:"
	@echo "  all             - Install dependencies and build the application (default)"
	@echo "  build           - Build the application"
	@echo "  build-demo      - Build the demo application"
	@echo "  build-all       - Build for all platforms"
	@echo "  build-windows   - Build for Windows"
	@echo "  build-linux     - Build for Linux"
	@echo "  build-macos     - Build for macOS"
	@echo "  run             - Run the application in normal mode"
	@echo "  run-api         - Run the application in API mode"
	@echo "  run-cli         - Run the application in CLI mode"
	@echo "  run-demo        - Run the demo application"
	@echo "  clean           - Clean build files"
	@echo "  deps            - Install dependencies"
	@echo "  test            - Run all tests"
	@echo "  test-protocols  - Run protocol tests"
	@echo "  test-managers   - Run manager tests"
	@echo "  docs            - Generate documentation"
	@echo "  help            - Show this help message"

# Release-related targets

# Version management
VERSION ?= dev

# Release target - prepares a new release
.PHONY: release
release:
	@echo "Preparing release $(VERSION)"
	./scripts/release.sh v$(VERSION)

# Build target - builds for all platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms"
	./scripts/build-all.sh $(VERSION)

# Package target - creates packages for all platforms
.PHONY: package
package:
	@echo "Creating packages"
	./scripts/package.sh $(VERSION)

# Clean release artifacts
.PHONY: clean-release
clean-release:
	rm -rf dist/ packages/ release/

# Help target
.PHONY: help-release
help-release:
	@echo "Release Management Commands:"
	@echo "  make release VERSION=X.Y.Z    - Prepare a new release"
	@echo "  make build-all VERSION=X.Y.Z  - Build for all platforms"
	@echo "  make package VERSION=X.Y.Z    - Create packages"
	@echo "  make clean-release            - Clean release artifacts"

.PHONY: all build build-demo build-all build-windows build-linux build-macos run run-api run-cli run-demo clean deps test test-protocols test-managers docs help