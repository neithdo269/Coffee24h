# ☕ Coffee24h — GDU Edition (Keep-Awake & System Power Daemon)

<p align="center">
  <img src="frontend/src/assets/images/coffee-logo.png" width="95" alt="Coffee24h GDU Blue Cup" />
  &nbsp;&nbsp;&nbsp;&nbsp;
  <img src="frontend/src/assets/images/coffee-leafcup.png" width="85" alt="Coffee24h GDU Eco Cup" />
</p>

<p align="center">
  <strong>High-Performance Keep-Awake Utility & Background Compute Persistence Daemon</strong><br />
  <em>Engineered for Gia Dinh University (GDU) — Conceived as a Core Component for a Future University OS</em>
</p>

<p align="center">
  <a href="https://github.com/nynwadhyhcz8/Coffee24h"><img src="https://img.shields.io/badge/Release-v1.0.0-blue.svg?style=flat-square" alt="Version" /></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8.svg?style=flat-square&logo=go" alt="Go Version" /></a>
  <a href="https://wails.io"><img src="https://img.shields.io/badge/Framework-Wails%20v2-DF1A29.svg?style=flat-square" alt="Wails v2" /></a>
  <a href="#-installation-guide-windows-macos-linux"><img src="https://img.shields.io/badge/Install-Windows%20%7C%20macOS%20%7C%20Linux-brightgreen.svg?style=flat-square" alt="Multi-OS Install" /></a>
  <a href="https://github.com/nynwadhyhcz8/Coffee24h/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License" /></a>
</p>

---

## 📑 Table of Contents
1. [🎓 About Gia Dinh University (GDU)](#-about-gia-dinh-university-gdu)
2. [🌌 The Mission: Paving the Way for a University Operating System](#-the-mission-paving-the-way-for-a-university-operating-system-gdu-os)
3. [⚡ Core Capabilities & Keep-Awake Modes](#-core-capabilities--keep-awake-modes)
4. [📥 Installation Guide (Windows, macOS, Linux)](#-installation-guide-windows-macos-linux)
   - [🪟 Windows Setup](#-1-windows-installation-windows-10--11)
   - [🍎 macOS Setup](#-2-macos-installation-apple-silicon--intel)
   - [🐧 Linux Setup](#-3-linux-installation-ubuntu-debian-fedora-arch)
5. [🖥️ Cross-Platform Architecture](#️-cross-platform-architecture)
6. [🛠️ Building from Source & Packaging](#️-building-from-source--packaging)
7. [⚙️ Various Ways to Set Up Coffee24h](#️-various-ways-to-set-up-coffee24h)
8. [💻 Command-Line Interface (CLI)](#-command-line-interface-cli)
9. [👨‍💻 Author & Mission Dedication](#-author--mission-dedication)

---

## 🎓 About Gia Dinh University (GDU)

**Trường Đại học Gia Định (Gia Dinh University - GDU)** is a progressive, forward-thinking university in Ho Chi Minh City, Vietnam. Rooted in the core philosophy of *"Ứng dụng - Trải nghiệm - Việc làm"* (Hands-on Application - Practical Experience - Employment Readiness) with the campus motto **"Chọn Gia Định - Chọn Thành Công"**, GDU places strong emphasis on real-world engineering, student innovation, and software craftsmanship.

### 🎨 Authentic GDU Visual Identity & Dual Editions
Coffee24h was crafted specifically around GDU's brand identity:
- **Colors:** Deep Oceanic Navy (`#061120`, `#0b2545`), Radiant Gold (`#f59e0b`), Flame Red (`#dc2626`).
- **👔 GDU Giảng Viên (Lecturer Edition):** Classical blue ceramic mug adorned with the official GDU emblem.
- **🌱 GDU Sinh Viên (Student Edition):** Eco-friendly leaf cup representing student dynamism, green campus sustainability, and fresh ideas.
- *Switchable dynamically on-the-fly with 1 click or CLI flags (`--student` / `--lecturer`).*

---

## 🌌 The Mission: Paving the Way for a University Operating System (GDU OS)

> *"Coffee24h was envisioned and coded by **Neith Do** on an ambitious mission: to research and contribute foundational, rock-solid system daemons toward the vision of a dedicated **University Operating System (GDU OS)** for Gia Dinh University."*

In modern university environments, research labs, and computer science departments, workstations are subjected to heavy, long-running computational workloads:
- **Machine Learning & Deep Neural Network Training** running over several hours.
- **Large Codebase Compilation, Kernel Building & 3D CAD Rendering**.
- **Automated Research Daemons & Lab Data Pipelines**.

When conventional operating systems arbitrarily enter standby or sleep mode, training checkpoints break, network sockets terminate, and critical student research is disrupted.

Coffee24h serves as a **production-grade proof-of-concept and subsystem module** for the future GDU OS, delivering:
1. **Low-level sleep interception** via direct kernel execution state registers (`SetThreadExecutionState` on Windows, `caffeinate` on macOS, `systemd-inhibit` on Linux).
2. **Pinned OS thread execution** (`runtime.LockOSThread()`) to prevent Go runtime goroutine migration from dropping OS wake flags.
3. **Smart hardware battery telemetry** to safeguard student laptops (<5% threshold auto-release).
4. **Shell auto-recovery** (`RegisterWindowMessageW("TaskbarCreated")`) ensuring tray persistence even if the desktop environment restarts.

---

## ⚡ Core Capabilities & Keep-Awake Modes

1. **Until Off Mode (Vô Hạn)**:
   - Blocks OS sleep, screen dimming, and screensaver.
   - Laptop stays active continuously until you stop it or the battery discharges completely.
2. **System-Only Awake Mode (Chạy Nền)**:
   - Allows screen to turn off to conserve battery and protect OLED/IPS displays, while keeping the CPU, background tasks, downloads, and rendering alive.
3. **Timed Mode (Hẹn Giờ)**:
   - Quick presets for **15m**, **30m**, **1h**, **2h**, and **4h** with active countdown timer and harmonic bell chimes.
4. **Smart Battery Safeguard (Bảo Vệ Pin)**:
   - Automatically releases awake mode if battery falls below 5% to prevent sudden data corruption and give the OS a chance to hibernate gracefully.
5. **Real-time Battery Telemetry**:
   - Detects AC status, battery percentage, and charging state dynamically.
6. **One-Click Workstation Lock (Win + L)**:
   - Instant security lock when stepping away from the desk in university libraries or labs.
7. **4-Corner Subpixel Antialiased Rounded Window**:
   - DirectComposition 32-bit ARGB with smooth 14px rounded corners and zero 1-bit GDI jagged pixel artifacts.

---

## 📥 Installation Guide (Windows, macOS, Linux)

### 🪟 1. Windows Installation (Windows 10 / 11)

#### Method A: NSIS Setup Wizard Installer
1. Download **`coffee24h-installer.exe`** from [Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Run the wizard: automatically installs Coffee24h, registers Desktop & Start Menu shortcuts, and provides clean uninstallation.

#### Method B: Standalone Portable Executable
1. Download **`coffee24h.exe`** from [Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Save it anywhere (e.g. `C:\Tools\Coffee24h\` or Desktop) and double-click to run. No admin rights or installation required.

#### Method C: 1-Click Multi-Setup Manager (From Repository)
Double-click **`setup.cmd`** in the repository root or run PowerShell:
```powershell
powershell -ExecutionPolicy Bypass -File .\setup.ps1 -All
```

---

### 🍎 2. macOS Installation (Apple Silicon & Intel)

Coffee24h is compiled as a Universal Binary supporting both **Apple Silicon (M1/M2/M3/M4)** and **Intel x86_64**.

#### Method A: Universal DMG Installer
1. Download **`Coffee24h-macOS-Universal.dmg`** from [Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Open the `.dmg` and drag **`Coffee24h.app`** into **`/Applications`**.
3. Launch from Launchpad or Spotlight.

> [!TIP]
> **macOS Gatekeeper Note:**  
> If prompted with *"Coffee24h cannot be opened because it is from an unidentified developer"*, right-click `Coffee24h.app` in Finder and choose **Open**, or run:
> ```bash
> xattr -cr /Applications/Coffee24h.app
> ```

#### Method B: Terminal Command-Line Alias
```bash
echo 'alias coffee24h="/Applications/Coffee24h.app/Contents/MacOS/Coffee24h"' >> ~/.zshrc
source ~/.zshrc
coffee24h --awake --time 2h
```

---

### 🐧 3. Linux Installation (Ubuntu, Debian, Fedora, Arch)

#### Method A: Debian / Ubuntu Package (`.deb`)
```bash
sudo apt install ./coffee24h_1.0.0_amd64.deb
# Or:
sudo dpkg -i coffee24h_1.0.0_amd64.deb
```
Installs binary to `/usr/bin/coffee24h`, configures system desktop entry, and registers high-res icons.

#### Method B: Universal Standalone AppImage
```bash
chmod +x Coffee24h-x86_64.AppImage
./Coffee24h-x86_64.AppImage
```

#### Method C: Manual Binary Installation
```bash
sudo install -m 755 coffee24h /usr/local/bin/coffee24h
sudo cp build/linux/coffee24h.desktop /usr/share/applications/
sudo cp build/linux/coffee24h.png /usr/share/icons/hicolor/256x256/apps/
```

---

## 🖥️ Cross-Platform Architecture

| OS | Keep-Awake Mechanism | Battery Detection | Workstation Lock |
| :--- | :--- | :--- | :--- |
| **Windows** | Native Win32 `SetThreadExecutionState` (`ES_CONTINUOUS \| ES_SYSTEM_REQUIRED \| ES_DISPLAY_REQUIRED`) | Win32 `GetSystemPowerStatus` | `user32.dll` `LockWorkStation` |
| **macOS** | Subprocess `caffeinate -i -d` | macOS `pmset -g batt` | `pmset displaysleepnow` |
| **Linux** | D-Bus / `systemd-inhibit` (fallback: `xdg-screensaver`) | Linux `/sys/class/power_supply/` | `loginctl lock-session` |

---

## 🛠️ Building from Source & Packaging

### Prerequisites
- **Go** (1.20+): `winget install GoLang.Go` (or `brew install go`, `sudo apt install golang`)
- **Node.js** (npm LTS): `winget install OpenJS.NodeJS.LTS` (or `brew install node`, `sudo apt install nodejs npm`)
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Linux Dev Headers:** `sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev` (or `webkit2gtk-4.1-dev`)

### 🪟 Windows Build
```powershell
# Standalone binary:
wails build

# NSIS Windows Setup Installer:
wails build -nsis
```

### 🍎 macOS Universal Build (.dmg)
```bash
chmod +x ./build/darwin/build-dmg.sh
./build/darwin/build-dmg.sh
```
Output: `build/bin/Coffee24h.app` & `build/bin/Coffee24h-macOS-Universal.dmg`

### 🐧 Linux Build (.deb & .AppImage)
```bash
# Debian / Ubuntu .deb package:
chmod +x ./build/linux/build-deb.sh
./build/linux/build-deb.sh 1.0.0

# Universal AppImage:
chmod +x ./build/linux/build-appimage.sh
./build/linux/build-appimage.sh
```
Output: `build/bin/coffee24h_1.0.0_amd64.deb` & `build/bin/Coffee24h-x86_64.AppImage`

---

## ⚙️ Various Ways to Set Up Coffee24h

Double-click **`setup.cmd`** in the root directory (or run `setup.ps1` via PowerShell) to open the interactive setup manager:

```text
==============================================================
   Coffee24h - Multi-Setup Manager (Various Setup Options)
   Keep-Awake Utility for Gia Dinh University (GDU)
==============================================================
  [1] Full Setup (All Recommended Setups in 1-Click)
  [2] Terminal CLI Setup (Add to PATH & WinGet Links)
  [3] Auto-Start on Windows Login (Silent to System Tray)
  [4] Create Desktop Shortcut (Custom GDU Icon)
  [5] Create Start Menu Shortcut (Windows Searchable)
  [6] Win + R Run Dialog ('coffee24h')
  [7] PowerShell Profile Aliases ('coffee', 'c24')
  [8] Remove / Uninstall All Integrations
==============================================================
```

---

## 💻 Command-Line Interface (CLI)

```text
Usage:
  coffee24h [flags]

Flags:
  -h, --help            Show CLI help reference
  -v, --version         Show application version
  -m, --tray            Start directly into System Tray ("Show hidden icons")
  -a, --awake           Immediately activate Keep-Awake mode upon launch
  -t, --time <dur>      Start a timed session (e.g. 15m, 30m, 1h, 2h, 4h, 90m)
  -s, --system-only     Keep CPU/downloads awake with display sleep enabled
  -l, --autolock        Auto-lock workstation (Win + L) when session expires
  --student             Launch in GDU Student Edition (Eco-friendly Leaf Cup)
  --lecturer            Launch in GDU Lecturer Edition (Blue GDU Logo Cup)
  --no-safeguard        Disable battery safeguard (<5% auto-stop)

Examples:
  coffee24h                    # Launch GUI normally
  coffee24h --tray             # Start silently in "Show hidden icons"
  coffee24h --awake            # Launch and immediately keep laptop awake
  coffee24h --time 2h          # Keep awake for 2 hours with timer
  coffee24h -s --time 1h       # 1 hour system-only awake (screen turns off, CPU runs)
  coffee24h --student          # Launch directly with GDU Student Edition
```

---

## 👨‍💻 Author & Mission Dedication

- **Author & Lead Developer:** **Neith Do** ([@nynwadhyhcz8](https://github.com/nynwadhyhcz8) / `16av112thien@gmail.com`)
- **Dedication:** Proudly designed and coded for **Gia Dinh University (GDU)**.
- **Mission:** A research stepping stone toward a future academic **University Operating System (GDU OS)**.
- **Architecture Documentation:** See [`DEVELOPER_GUIDE.md`](../DEVELOPER_GUIDE.md) for full deep-dive architectural notes.
