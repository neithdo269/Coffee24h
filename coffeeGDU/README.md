# Coffee24h ☕
# CoffeeGDU is what Vietnamese students call Coffee24h!
A lightweight, modern cross-platform utility that keeps your laptop awake until power runs out, until a timer expires, or for background processing.
A special ordered by Neith Do using vibe coding @2026.
## 🚀 Key Features

1. **Until Off Mode (Infinite Awake)**:
   - Blocks OS sleep, screen dimming, and screensaver.
   - Laptop stays active continuously until you stop it or the battery discharges completely.
2. **System-Only Awake Mode**:
   - Allows screen to turn off to conserve battery and protect displays (e.g., OLED), while keeping the CPU, background tasks, downloads, and rendering alive.
3. **Timed Mode**:
   - Quick presets for 15m, 30m, 1h, 2h, and 4h with active countdown timer.
4. **Smart Battery Safeguard**:
   - Automatically releases awake mode if battery falls below 5% to prevent sudden data corruption and give the OS a chance to hibernate gracefully (can be toggled on/off).
5. **Real-time Battery Telemetry**:
   - Detects AC status, battery percentage, and charging state.

---

## 🖥️ Cross-Platform Architecture

| OS | Keep-Awake Mechanism | Battery Detection |
| :--- | :--- | :--- |
| **Windows** | Native Win32 `SetThreadExecutionState` (`ES_CONTINUOUS \| ES_SYSTEM_REQUIRED \| ES_DISPLAY_REQUIRED`) | Win32 `GetSystemPowerStatus` |
| **macOS** | Subprocess `caffeinate -i -d` | macOS `pmset -g batt` |
| **Linux** | D-Bus / `systemd-inhibit` (fallback: `xdg-screensaver`) | Linux `/sys/class/power_supply/` |

---

## 🛠️ Building & Running

> [!IMPORTANT]
> This application is built with **Wails v2**. If you compile with plain `go build`, the application will display a popup: *"Wails applications will not build without the correct build tags."*
> You must build using either the **Wails CLI** or `go build` with the required tags.

### Prerequisites
- **Go** (1.20+): `winget install GoLang.Go`
- **Node.js** (npm): `winget install OpenJS.NodeJS.LTS`
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **WebView2 Runtime**: Preinstalled on Windows 10/11

---

### Option 1: Standard Wails Build (Recommended)
From the project folder (`coffeeGDU`):
```powershell
# Compiles frontend with Vite and packages the executable into build/bin/
wails build
```

### Option 2: Manual Go Build
If compiling directly with `go build`, build the frontend first and pass the required `-tags` and `-ldflags`:
```powershell
# 1. Build frontend
cd frontend
npm install
npm run build
cd ..

# 2. Build Go binary with Wails tags
go build -tags desktop,production -ldflags "-H windowsgui -w -s" -o coffee24h.exe .
```

---

## 📦 Cross-Platform Packaging

### 🍎 macOS (.app / .dmg)
To build a universal macOS application (supporting both Intel & Apple Silicon M1/M2/M3/M4) and package into a `.dmg` installer:
```bash
# 1. Standard build
wails build -platform darwin/universal

# 2. Or package into a clean .dmg installer:
chmod +x ./build/darwin/build-dmg.sh
./build/darwin/build-dmg.sh
```
Output: `build/bin/Coffee24h.app` & `build/bin/Coffee24h-macOS-Universal.dmg`

### 🐧 Linux (.deb / AppImage)
To compile for Linux and create official Debian/Ubuntu `.deb` or standalone `.AppImage` packages:
```bash
# 1. Standard Linux build
wails build -platform linux/amd64

# 2. Build Debian/Ubuntu (.deb) installer:
chmod +x ./build/linux/build-deb.sh
./build/linux/build-deb.sh 1.0.0

# 3. Build standalone AppImage:
chmod +x ./build/linux/build-appimage.sh
./build/linux/build-appimage.sh
```
Output: `build/bin/coffee24h_1.0.0_amd64.deb` & `build/bin/Coffee24h-x86_64.AppImage`

### 🪟 Windows (.exe / NSIS Installer)
```powershell
# Standalone executable:
wails build

# Or generate full Windows Setup Installer (NSIS):
wails build -nsis
```
Output: `build/bin/coffee24h.exe` & `build/bin/coffee24h-installer.exe`

---

### Option 3: Quick Browser UI Preview (Zero Go Dependencies)
To preview and test the UI layout without compiling the Go desktop binary:
```powershell
powershell -ExecutionPolicy Bypass -File .\serve.ps1
# Open http://localhost:8080 in your browser
```
