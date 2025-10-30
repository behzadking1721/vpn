#!/bin/bash

# Release script for VPN client - manages versioning and GitHub releases

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display usage
usage() {
    echo "Usage: $0 <version> [options]"
    echo "  version: Version to release (e.g., v1.0.0)"
    echo ""
    echo "Options:"
    echo "  --no-build     Skip building binaries"
    echo "  --no-package   Skip packaging"
    echo "  --dry-run      Show what would be done without doing it"
    echo "  --help         Display this help message"
    exit 1
}

# Parse arguments
if [[ $# -eq 0 ]] || [[ "$1" == "--help" ]]; then
    usage
fi

VERSION="$1"
shift

# Default options
BUILD=true
PACKAGE=true
DRY_RUN=false

# Parse options
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-build)
            BUILD=false
            shift
            ;;
        --no-package)
            PACKAGE=false
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            ;;
    esac
done

# Validate version format (should be vX.Y.Z)
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Version must be in format vX.Y.Z (e.g., v1.0.0)${NC}"
    exit 1
fi

echo -e "${BLUE}Preparing to release VPN Client ${VERSION}${NC}"

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${YELLOW}Warning: Working directory is not clean${NC}"
    if [ "$DRY_RUN" = false ]; then
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi

# Update version in versioninfo.json
if [ -f "versioninfo.json" ]; then
    echo -e "${YELLOW}Updating version in versioninfo.json...${NC}"
    if [ "$DRY_RUN" = false ]; then
        # Extract version numbers
        MAJOR=$(echo $VERSION | cut -d. -f1 | sed 's/v//')
        MINOR=$(echo $VERSION | cut -d. -f2)
        PATCH=$(echo $VERSION | cut -d. -f3)
        
        # Update JSON using sed
        sed -i.bak \
            -e "s/\"Major\": [0-9]*,/\"Major\": $MAJOR,/" \
            -e "s/\"Minor\": [0-9]*,/\"Minor\": $MINOR,/" \
            -e "s/\"Patch\": [0-9]*,/\"Patch\": $PATCH,/" \
            -e "s/\"FileVersion\": \"[0-9.]*\"/\"FileVersion\": \"$MAJOR.$MINOR.$PATCH.0\"/" \
            -e "s/\"ProductVersion\": \"[0-9.]*\"/\"ProductVersion\": \"$MAJOR.$MINOR.$PATCH.0\"/" \
            versioninfo.json
        
        rm versioninfo.json.bak
    fi
    echo -e "${GREEN}✓ Updated versioninfo.json${NC}"
fi

# Update version in NSIS script
if [ -f "installer.nsi" ]; then
    echo -e "${YELLOW}Updating version in installer.nsi...${NC}"
    if [ "$DRY_RUN" = false ]; then
        sed -i.bak "s/!define APP_VERSION \".*\"/!define APP_VERSION \"${VERSION#v}\"/" installer.nsi
        rm installer.nsi.bak
    fi
    echo -e "${GREEN}✓ Updated installer.nsi${NC}"
fi

# Build binaries if requested
if [ "$BUILD" = true ]; then
    echo -e "${YELLOW}Building binaries...${NC}"
    if [ "$DRY_RUN" = false ]; then
        ./scripts/build-all.sh "${VERSION#v}"
    else
        echo -e "${BLUE}[DRY RUN] Would run: ./scripts/build-all.sh ${VERSION#v}${NC}"
    fi
else
    echo -e "${YELLOW}Skipping build phase${NC}"
fi

# Package if requested
if [ "$PACKAGE" = true ]; then
    echo -e "${YELLOW}Creating packages...${NC}"
    if [ "$DRY_RUN" = false ]; then
        ./scripts/package.sh "${VERSION#v}"
    else
        echo -e "${BLUE}[DRY RUN] Would run: ./scripts/package.sh ${VERSION#v}${NC}"
    fi
else
    echo -e "${YELLOW}Skipping packaging phase${NC}"
fi

# Create Git tag
echo -e "${YELLOW}Creating Git tag...${NC}"
if [ "$DRY_RUN" = false ]; then
    git tag -a "$VERSION" -m "Release $VERSION"
else
    echo -e "${BLUE}[DRY RUN] Would create Git tag: $VERSION${NC}"
fi
echo -e "${GREEN}✓ Tag created${NC}"

# Summary
echo -e "${GREEN}"
echo "=================================================="
echo "           RELEASE PREPARATION COMPLETE"
echo "=================================================="
echo "Version:          ${VERSION}"
echo "Build performed:  ${BUILD}"
echo "Packaging done:   ${PACKAGE}"
echo "Git tag created:  Yes"
echo ""
echo "Next steps:"
echo "1. Review the build artifacts in dist/ and packages/"
echo "2. Test the packages on target platforms"
echo "3. Push the tag: git push origin ${VERSION}"
echo "4. Create GitHub release: gh release create ${VERSION} packages/*"
echo "=================================================="
echo -e "${NC}"