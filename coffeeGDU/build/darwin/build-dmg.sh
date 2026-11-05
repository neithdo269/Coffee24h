#!/usr/bin/env bash
# ========================================================
# Script: Package Coffee24h into a macOS .dmg Installer
# ========================================================
set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUILD_DIR="${ROOT_DIR}/build/bin"
APP_NAME="Coffee24h"
DMG_NAME="Coffee24h-macOS-Universal.dmg"

echo "☕ Building Coffee24h macOS Universal .app..."

cd "${ROOT_DIR}"
# Build universal binary (both Intel x86_64 and Apple Silicon arm64)
wails build -platform darwin/universal -clean

if [ ! -d "${BUILD_DIR}/${APP_NAME}.app" ]; then
    echo "❌ Error: ${APP_NAME}.app not found in build/bin/!"
    exit 1
fi

echo "📦 Packaging into DMG..."
rm -f "${BUILD_DIR}/${DMG_NAME}"

# Create DMG with Applications shortcut
DMG_TMP="/tmp/${APP_NAME}-dmg-staging"
rm -rf "${DMG_TMP}"
mkdir -p "${DMG_TMP}"

cp -R "${BUILD_DIR}/${APP_NAME}.app" "${DMG_TMP}/"
ln -s /Applications "${DMG_TMP}/Applications"

hdiutil create -volname "${APP_NAME}" -srcfolder "${DMG_TMP}" -ov -format UDZO "${BUILD_DIR}/${DMG_NAME}"
rm -rf "${DMG_TMP}"

echo "✅ Successfully created macOS installer: ${BUILD_DIR}/${DMG_NAME}"
