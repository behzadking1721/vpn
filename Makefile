# Makefile for VPN Client

# Variables
VERSION ?= dev
BUILD_TIME := $(shell date +%Y%m%d%H%M%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
DIRTY := $(shell git diff --quiet || echo "-dirty")
GIT_COMMIT := $(GIT_COMMIT)$(DIRTY)

# Build flags
LDFLAGS := -ldflags="-X 'main.version=$(VERSION)' -X 'main.buildTime=$(BUILD_TIME)' -X 'main.gitCommit=$(GIT_COMMIT)'"

# Platform-specific variables
WINDOWS_OUT := dist/windows/vpn-client.exe
LINUX_OUT := dist/linux/vpn-client
MACOS_OUT := dist/macos/vpn-client

# Default target
.PHONY: all
all: clean build-all

# Clean build directory
.PHONY: clean
clean:
	rm -rf dist/

# Build for all platforms
.PHONY: build-all
build-all: build-windows build-linux build-macos

# Build for Windows
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(WINDOWS_OUT) ./cmd/vpn-client
	@echo "Windows build completed: $(WINDOWS_OUT)"

# Build for Linux
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(LINUX_OUT) ./cmd/vpn-client
	@echo "Linux build completed: $(LINUX_OUT)"

# Build for macOS
.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(MACOS_OUT) ./cmd/vpn-client
	@echo "macOS build completed: $(MACOS_OUT)"

# Build CLI tool for all platforms
.PHONY: build-cli-all
build-cli-all: build-cli-windows build-cli-linux build-cli-macos

# Build CLI for Windows
.PHONY: build-cli-windows
build-cli-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/windows/vpnctl.exe ./cmd/vpnctl
	@echo "Windows CLI build completed: dist/windows/vpnctl.exe"

# Build CLI for Linux
.PHONY: build-cli-linux
build-cli-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/linux/vpnctl ./cmd/vpnctl
	@echo "Linux CLI build completed: dist/linux/vpnctl"

# Build CLI for macOS
.PHONY: build-cli-macos
build-cli-macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/macos/vpnctl ./cmd/vpnctl
	@echo "macOS CLI build completed: dist/macos/vpnctl"

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	@echo "Dependencies updated"

# Create release directory structure
.PHONY: release-structure
release-structure:
	mkdir -p releases/$(VERSION)/windows
	mkdir -p releases/$(VERSION)/linux
	mkdir -p releases/$(VERSION)/macos
	mkdir -p releases/$(VERSION)/mobile/android
	mkdir -p releases/$(VERSION)/mobile/ios
	@echo "Release directory structure created for version $(VERSION)"

# Help target
.PHONY: help
help:
	@echo "VPN Client Makefile"
	@echo "=================="
	@echo "Available targets:"
	@echo "  all               - Clean and build for all platforms"
	@echo "  clean             - Remove build artifacts"
	@echo "  build-all         - Build for all platforms"
	@echo "  build-windows     - Build for Windows"
	@echo "  build-linux       - Build for Linux"
	@echo "  build-macos       - Build for macOS"
	@echo "  build-cli-all     - Build CLI for all platforms"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  deps              - Install dependencies"
	@echo "  release-structure - Create release directory structure"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION           - Release version (default: dev)"