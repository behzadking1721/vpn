#!/bin/bash

# Packaging script for VPN client - creates installers for various platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get version from command line or git
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}

# Base directories
DIST_DIR="dist"
PACKAGE_DIR="packages"
RELEASE_DIR="release"

# Clean previous packages
echo -e "${BLUE}Cleaning previous packages...${NC}"
rm -rf ${PACKAGE_DIR} ${RELEASE_DIR}
mkdir -p ${PACKAGE_DIR} ${RELEASE_DIR}

echo -e "${BLUE}Packaging VPN Client v${VERSION}${NC}"

# Check if builds exist
if [ ! -d "${DIST_DIR}" ]; then
    echo -e "${RED}Error: Build directory '${DIST_DIR}' does not exist. Run build script first.${NC}"
    exit 1
fi

# Package for Windows (NSIS installer)
if command -v makensis &> /dev/null; then
    echo -e "${YELLOW}Creating Windows installer...${NC}"
    
    # Update version in NSIS script
    sed -i.bak "s/!define APP_VERSION \".*\"/!define APP_VERSION \"${VERSION}\"/" installer.nsi
    
    # Build installer
    makensis installer.nsi
    
    if [ -f "VPN Client-installer.exe" ]; then
        mv "VPN Client-installer.exe" "${PACKAGE_DIR}/vpn-client-windows-${VERSION}.exe"
        echo -e "${GREEN}✓ Windows installer created${NC}"
    else
        echo -e "${RED}Failed to create Windows installer${NC}"
    fi
    
    # Restore NSIS script
    mv installer.nsi.bak installer.nsi
else
    echo -e "${YELLOW}NSIS not found, skipping Windows installer creation${NC}"
fi

# Create ZIP archives for all platforms
echo -e "${YELLOW}Creating ZIP archives...${NC}"

for platform_dir in "${DIST_DIR}"/*; do
    if [ -d "$platform_dir" ]; then
        platform=$(basename "$platform_dir")
        zip_name="vpn-client-${platform}-${VERSION}.zip"
        
        echo -e "${BLUE}Creating archive for ${platform}...${NC}"
        cd "${DIST_DIR}"
        zip -r "../${PACKAGE_DIR}/${zip_name}" "${platform}" > /dev/null
        cd ..
        echo -e "${GREEN}✓ Created ${zip_name}${NC}"
    fi
done

# Create tar.gz archives for Unix-like systems
echo -e "${YELLOW}Creating tar.gz archives...${NC}"

for platform_dir in "${DIST_DIR}"/*; do
    if [ -d "$platform_dir" ]; then
        platform=$(basename "$platform_dir")
        
        # Skip Windows platforms for tar.gz
        if [[ "$platform" != windows* ]]; then
            tar_name="vpn-client-${platform}-${VERSION}.tar.gz"
            
            echo -e "${BLUE}Creating tarball for ${platform}...${NC}"
            tar -czf "${PACKAGE_DIR}/${tar_name}" -C "${DIST_DIR}" "${platform}" > /dev/null
            echo -e "${GREEN}✓ Created ${tar_name}${NC}"
        fi
    fi
done

# Create Debian package for Linux
if command -v dpkg-deb &> /dev/null; then
    echo -e "${YELLOW}Creating Debian package...${NC}"
    
    # Create directory structure for DEB package
    DEB_DIR="${PACKAGE_DIR}/deb"
    mkdir -p "${DEB_DIR}/usr/bin"
    mkdir -p "${DEB_DIR}/usr/share/applications"
    mkdir -p "${DEB_DIR}/usr/share/icons/hicolor/256x256/apps"
    mkdir -p "${DEB_DIR}/DEBIAN"
    
    # Copy main executable (use amd64 as example)
    if [ -f "${DIST_DIR}/linux_amd64/vpn-client" ]; then
        cp "${DIST_DIR}/linux_amd64/vpn-client" "${DEB_DIR}/usr/bin/"
        chmod +x "${DEB_DIR}/usr/bin/vpn-client"
    fi
    
    # Create control file
    cat > "${DEB_DIR}/DEBIAN/control" << EOF
Package: vpn-client
Version: ${VERSION}
Section: net
Priority: optional
Architecture: amd64
Maintainer: Your Name <your.email@example.com>
Description: Multi-Protocol VPN Client
 A feature-rich VPN client supporting multiple protocols
EOF
    
    # Create desktop entry
    cat > "${DEB_DIR}/usr/share/applications/vpn-client.desktop" << EOF
[Desktop Entry]
Name=VPN Client
Comment=Multi-Protocol VPN Client
Exec=/usr/bin/vpn-client
Icon=vpn-client
Terminal=false
Type=Application
Categories=Network;VPN;
EOF
    
    # Build DEB package
    dpkg-deb --build "${DEB_DIR}" "${PACKAGE_DIR}/vpn-client_${VERSION}_amd64.deb" > /dev/null
    rm -rf "${DEB_DIR}"
    echo -e "${GREEN}✓ Debian package created${NC}"
else
    echo -e "${YELLOW}dpkg-deb not found, skipping Debian package creation${NC}"
fi

echo -e "${GREEN}All packaging completed successfully!${NC}"
echo -e "${BLUE}Packages located in: ${PACKAGE_DIR}${NC}"