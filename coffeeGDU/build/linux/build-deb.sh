#!/usr/bin/env bash
# ========================================================
# Script: Build Debian/Ubuntu (.deb) package for Coffee24h
# Usage: ./build-deb.sh [version]
# ========================================================
set -e

VERSION="${1:-1.0.0}"
ARCH="amd64"
PKG_NAME="coffee24h_${VERSION}_${ARCH}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUILD_DIR="${ROOT_DIR}/build/bin"
PKG_DIR="${ROOT_DIR}/build/linux/${PKG_NAME}"

echo "☕ Building Coffee24h Debian package v${VERSION}..."

# Ensure binary exists
if [ ! -f "${BUILD_DIR}/coffee24h" ]; then
    echo "⚙️ Binary not found in build/bin/. Compiling with Wails..."
    cd "${ROOT_DIR}"
    wails build -platform linux/amd64
fi

# Clean previous package directory
rm -rf "${PKG_DIR}" "${PKG_DIR}.deb"

# Create Debian file structure
mkdir -p "${PKG_DIR}/DEBIAN"
mkdir -p "${PKG_DIR}/usr/bin"
mkdir -p "${PKG_DIR}/usr/share/applications"
mkdir -p "${PKG_DIR}/usr/share/icons/hicolor/256x256/apps"

# 1. Copy Binary
cp "${BUILD_DIR}/coffee24h" "${PKG_DIR}/usr/bin/coffee24h"
chmod +x "${PKG_DIR}/usr/bin/coffee24h"

# 2. Copy Desktop Entry
cp "${ROOT_DIR}/build/linux/coffee24h.desktop" "${PKG_DIR}/usr/share/applications/coffee24h.desktop"

# 3. Copy Icon
cp "${ROOT_DIR}/build/linux/coffee24h.png" "${PKG_DIR}/usr/share/icons/hicolor/256x256/apps/coffee24h.png"

# 4. Generate Control file
cat <<EOF > "${PKG_DIR}/DEBIAN/control"
Package: coffee24h
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Neith Do <16av112thien@gmail.com>
Depends: libgtk-3-0, libwebkit2gtk-4.0-37 | libwebkit2gtk-4.1-0
Description: Coffee24h - Cross-Platform Keep-Awake Utility
 A lightweight, modern cross-platform utility that keeps your laptop awake
 until power runs out, until a timer expires, or for background processing.
 Part of the GDU (Gia Dinh University) Window Manager Project.
EOF

# 5. Build .deb package
dpkg-deb --build --root-owner-group "${PKG_DIR}" "${BUILD_DIR}/${PKG_NAME}.deb"

# Clean temporary package directory
rm -rf "${PKG_DIR}"

echo "✅ Successfully built: ${BUILD_DIR}/${PKG_NAME}.deb"
EOF