#!/usr/bin/env bash
# ========================================================
# Script: Build Linux AppImage for Coffee24h
# ========================================================
set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUILD_DIR="${ROOT_DIR}/build/bin"
APP_DIR="${ROOT_DIR}/build/linux/AppDir"

echo "☕ Building Coffee24h AppImage..."

if [ ! -f "${BUILD_DIR}/coffee24h" ]; then
    echo "⚙️ Compiling with Wails..."
    cd "${ROOT_DIR}"
    wails build -platform linux/amd64
fi

rm -rf "${APP_DIR}"
mkdir -p "${APP_DIR}/usr/bin"
mkdir -p "${APP_DIR}/usr/share/icons/hicolor/256x256/apps"

cp "${BUILD_DIR}/coffee24h" "${APP_DIR}/usr/bin/coffee24h"
chmod +x "${APP_DIR}/usr/bin/coffee24h"

cp "${ROOT_DIR}/build/linux/coffee24h.desktop" "${APP_DIR}/coffee24h.desktop"
cp "${ROOT_DIR}/build/linux/coffee24h.png" "${APP_DIR}/coffee24h.png"
cp "${ROOT_DIR}/build/linux/coffee24h.png" "${APP_DIR}/usr/share/icons/hicolor/256x256/apps/coffee24h.png"

# Create AppRun entrypoint
cat <<'EOF' > "${APP_DIR}/AppRun"
#!/bin/sh
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
export PATH="${HERE}/usr/bin:${PATH}"
export LD_LIBRARY_PATH="${HERE}/usr/lib:${LD_LIBRARY_PATH}"
exec "${HERE}/usr/bin/coffee24h" "$@"
EOF
chmod +x "${APP_DIR}/AppRun"

# Download appimagetool if not present and build
if ! command -v appimagetool &> /dev/null; then
    echo "📥 Downloading appimagetool..."
    curl -sL -o /tmp/appimagetool "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage"
    chmod +x /tmp/appimagetool
    TOOL="/tmp/appimagetool"
else
    TOOL="appimagetool"
fi

ARCH=x86_64 "$TOOL" "${APP_DIR}" "${BUILD_DIR}/Coffee24h-x86_64.AppImage"
rm -rf "${APP_DIR}"
echo "✅ Successfully built: ${BUILD_DIR}/Coffee24h-x86_64.AppImage"
