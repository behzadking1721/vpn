#!/bin/bash

# Release script for VPN Client
# Usage: ./scripts/release.sh [version]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[STATUS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if git is clean
check_git_status() {
    if [[ -n $(git status --porcelain) ]]; then
        print_error "Git working directory is not clean. Please commit or stash changes."
        exit 1
    fi
}

# Get current version
get_current_version() {
    # Try to get version from versioninfo.json
    if [ -f "versioninfo.json" ]; then
        grep -o '"FileVersion": "[0-9.]*"' versioninfo.json | cut -d'"' -f4
    else
        echo "0.0.0"
    fi
}

# Update version in files
update_version() {
    local version=$1
    local current_version=$2
    
    print_status "Updating version from $current_version to $version"
    
    # Update versioninfo.json
    if [ -f "versioninfo.json" ]; then
        sed -i.bak "s/\"FileVersion\": \"$current_version\"/\"FileVersion\": \"$version\"/" versioninfo.json
        sed -i.bak "s/\"ProductVersion\": \"$current_version\"/\"ProductVersion\": \"$version\"/" versioninfo.json
        rm versioninfo.json.bak
    fi
    
    # Update installer.nsi
    if [ -f "installer.nsi" ]; then
        sed -i.bak "s/!define APP_VERSION \"$current_version\"/!define APP_VERSION \"$version\"/" installer.nsi
        rm installer.nsi.bak
    fi
    
    # Update changelog
    if [ -f "CHANGELOG.md" ]; then
        # Add new version to changelog
        sed -i.bak "1s/^/## [$version] - $(date +%Y-%m-%d)\n\n### Added\n-\n\n### Changed\n-\n\n### Fixed\n-\n\n/" CHANGELOG.md
        rm CHANGELOG.md.bak
    else
        # Create changelog
        cat > CHANGELOG.md << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [$version] - $(date +%Y-%m-%d)

### Added
-

### Changed
-

### Fixed
-
EOF
    fi
}

# Create git tag
create_tag() {
    local version=$1
    print_status "Creating git tag v$version"
    git add .
    git commit -m "Release version $version"
    git tag -a "v$version" -m "Release version $version"
}

# Build application
build_application() {
    local version=$1
    print_status "Building application version $version"
    
    # Clean previous builds
    rm -rf build dist bin
    
    # Build using make
    make build-all
    
    # Create distribution packages
    print_status "Creating distribution packages"
    
    # This would be where you'd create your actual installers
    # For now, we'll just create zip archives
    mkdir -p releases/v$version
    
    # Copy built binaries to releases directory
    if [ -d "bin" ]; then
        cp -r bin/* releases/v$version/
    fi
    
    # Create a simple archive for each platform
    for dir in dist/*; do
        if [ -d "$dir" ]; then
            platform=$(basename "$dir")
            if [ "$platform" != "*" ]; then
                zip -r "releases/v$version/vpn-client-$version-$platform.zip" "$dir"
            fi
        fi
    done
}

# Main function
main() {
    local version=$1
    
    if [ -z "$version" ]; then
        print_error "Version not specified"
        echo "Usage: $0 <version>"
        echo "Example: $0 1.2.0"
        exit 1
    fi
    
    # Validate version format (semantic versioning)
    if ! [[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "Invalid version format. Use semantic versioning (e.g., 1.2.0)"
        exit 1
    fi
    
    print_status "Starting release process for version $version"
    
    # Check git status
    check_git_status
    
    # Get current version
    current_version=$(get_current_version)
    print_status "Current version: $current_version"
    
    # Update version in files
    update_version "$version" "$current_version"
    
    # Build application
    build_application "$version"
    
    # Create git tag
    create_tag "$version"
    
    print_status "Release v$version created successfully!"
    print_status "Release files are located in releases/v$version/"
    print_status "Don't forget to push the changes and tag:"
    echo "  git push origin main"
    echo "  git push origin v$version"
}

# Run main function with all arguments
main "$@"