#!/bin/bash

# Enhanced build script for VPN client supporting multiple platforms and architectures

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the version from command line argument or git tag
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Output directory
OUTPUT_BASE_DIR="dist"
RELEASE_DIR="release"

# Clean previous builds
echo -e "${BLUE}Cleaning previous builds...${NC}"
rm -rf ${OUTPUT_BASE_DIR} ${RELEASE_DIR}
mkdir -p ${OUTPUT_BASE_DIR} ${RELEASE_DIR}

echo -e "${BLUE}Building VPN Client v${VERSION} (${COMMIT_HASH})${NC}"

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH}"

# Platform configurations
declare -A PLATFORMS=(
    ["windows/amd64"]=""
    ["windows/386"]=""
    ["windows/arm64"]=""
    ["linux/amd64"]=""
    ["linux/386"]=""
    ["linux/arm64"]=""
    ["linux/arm"]=""
    ["darwin/amd64"]=""
    ["darwin/arm64"]=""
)

# Build for each platform
for platform in "${!PLATFORMS[@]}"; do
    IFS='/' read -ra ADDR <<< "$platform"
    GOOS=${ADDR[0]}
    GOARCH=${ADDR[1]}
    
    # Determine binary name
    BIN_NAME="vpn-client"
    if [ "$GOOS" = "windows" ]; then
        BIN_NAME="${BIN_NAME}.exe"
    fi
    
    # Special handling for macOS
    if [ "$GOOS" = "darwin" ]; then
        # Enable CGO for macOS
        export CGO_ENABLED=1
    else
        export CGO_ENABLED=0
    fi
    
    OUTPUT_DIR="${OUTPUT_BASE_DIR}/${GOOS}_${GOARCH}"
    mkdir -p "${OUTPUT_DIR}"
    
    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Build command
    env GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${OUTPUT_DIR}/${BIN_NAME}" ./src
    
    # Copy UI files for non-webview platforms
    if [ -d "ui/desktop" ]; then
        cp -r ui/desktop "${OUTPUT_DIR}/ui"
    fi
    
    # Create basic directories
    mkdir -p "${OUTPUT_DIR}/config"
    mkdir -p "${OUTPUT_DIR}/data"
    mkdir -p "${OUTPUT_DIR}/logs"
    
    echo -e "${GREEN}✓ Built for ${GOOS}/${GOARCH}${NC}"
done

echo -e "${GREEN}All builds completed successfully!${NC}"
echo -e "${BLUE}Output directory: ${OUTPUT_BASE_DIR}${NC}"