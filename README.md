# ☕ Coffee24h — GDU Keep-Awake System Utility

<p align="center">
  <img src="coffeeGDU/frontend/src/assets/images/coffee-logo.png" width="100" alt="Coffee24h GDU Blue Cup" />
  &nbsp;&nbsp;&nbsp;&nbsp;
  <img src="coffeeGDU/frontend/src/assets/images/coffee-leafcup.png" width="90" alt="Coffee24h GDU Eco Cup" />
</p>

<p align="center">
  <strong>Cross-Platform Keep-Awake Daemon & Intelligent Power Orchestration Utility</strong><br />
  <em>Specially crafted for Gia Dinh University (GDU) Faculty & Students</em>
</p>

<p align="center">
  <a href="https://github.com/nynwadhyhcz8/Coffee24h"><img src="https://img.shields.io/badge/Release-v1.0.0-blue.svg?style=flat-square" alt="Version" /></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8.svg?style=flat-square&logo=go" alt="Go Version" /></a>
  <a href="https://wails.io"><img src="https://img.shields.io/badge/Framework-Wails%20v2-DF1A29.svg?style=flat-square" alt="Wails v2" /></a>
  <a href="#-installation-guide-windows-macos-linux-npm"><img src="https://img.shields.io/badge/Install-Windows%20%7C%20macOS%20%7C%20Linux%20%7C%20NPM-brightgreen.svg?style=flat-square" alt="Multi-OS Install" /></a>
  <a href="https://nodejs.org"><img src="https://img.shields.io/badge/NPM-Node.js%20Ready-CB3837.svg?style=flat-square&logo=npm" alt="NPM" /></a>
  <a href="#-about-gia-dinh-university-gdu--brand-identity"><img src="https://img.shields.io/badge/Edition-GDU%20University-f59e0b.svg?style=flat-square" alt="GDU Edition" /></a>
  <a href="https://github.com/nynwadhyhcz8/Coffee24h/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License" /></a>
</p>

---

## 📑 Table of Contents
1. [🎓 About Gia Dinh University (GDU) & Brand Identity](#-about-gia-dinh-university-gdu--brand-identity)
2. [🌌 The Vision & Mission: Contributing to a University Operating System](#-the-vision--mission-contributing-to-a-university-operating-system)
3. [⚡ Core Capabilities & Keep-Awake Modes](#-core-capabilities--keep-awake-modes)
4. [📥 Installation Guide (Windows, macOS, Linux, NPM)](#-installation-guide-windows-macos-linux-npm)
   - [🪟 1. Windows Installation](#-1-windows-installation-windows-10--11)
   - [🍎 2. macOS Installation](#-2-macos-installation-apple-silicon--intel)
   - [🐧 3. Linux Installation](#-3-linux-installation-ubuntu-debian-fedora-arch)
   - [📦 4. NPM & Node.js (CLI & Web Preview)](#-4-npm--nodejs-cli--web-preview)
5. [🏗️ Low-Level OS Architecture](#️-low-level-os-architecture)
6. [🖥️ Cross-Platform Hardware & Power Mapping](#️-cross-platform-hardware--power-mapping)
7. [🛠️ Build from Source & Packaging](#️-build-from-source--packaging)
8. [⚙️ 1-Click Multi-Setup Manager](#️-1-click-multi-setup-manager)
9. [💻 Command-Line Interface (CLI)](#-command-line-interface-cli)
10. [👨‍💻 Author & Acknowledgments](#-author--acknowledgments)

---

## 🎓 About Gia Dinh University (GDU) & Brand Identity

**Trường Đại học Gia Định (Gia Dinh University - GDU)** is a dynamic, modern, and forward-looking university based in Ho Chi Minh City, Vietnam. Guided by the educational philosophy of *"Ứng dụng - Trải nghiệm - Việc làm"* (Application - Practical Experience - Career Readiness) and the motto **"Chọn Gia Định - Chọn Thành Công"** (Choose Gia Dinh - Choose Success), GDU empowers students to master practical, industry-grade software engineering, computer science, and creative technology.

<p align="center">
  <img src="coffeeGDU/frontend/src/assets/images/logo-universal.png" width="300" alt="GDU Universal Logo" />
</p>

### 🎨 Authentic GDU Visual Design
Coffee24h was engineered from the ground up to reflect the proud visual identity of Gia Dinh University:
- **GDU Navy & Deep Ocean (`#061120`, `#0b2545`):** Expressing academic depth, technological trust, and focus.
- **GDU Gold (`#f59e0b`):** Radiating energy, creativity, and student success.
- **Flame Red (`#dc2626`):** Symbolizing vibrant enthusiasm and determination.

### 🔄 Dual-Edition Dynamic Toggle (Clean Code Architecture)
Users can seamlessly toggle between two distinct university editions with a single click on the cup or title tag, or by passing CLI flags:
- **👔 GDU Giảng Viên (Lecturer Edition):** Features the distinguished Royal Blue ceramic cup stamped with the official GDU emblem.
- **🌱 GDU Sinh Viên (Student Edition):** Features the Eco-friendly Recyclable Leaf Cup, representing youthful vitality, green campus initiatives, and sustainable lifestyle.

---

## 🌌 The Vision & Mission: Contributing to a University Operating System

### 💡 The Backstory & Passion
> *"Coffee24h was designed and coded by **Neith Do** on a dedicated technological mission: to research, engineer, and contribute foundational low-level system modules toward the long-term vision of a specialized **University Operating System (GDU OS)** tailored for Gia Dinh University."*

Creating a custom or specialized operating system distribution for academic institutions is a profound endeavor. In university research centers, computer labs, and lecture halls, students and professors regularly execute long-running computational workloads:
- **Machine Learning & Deep Neural Network Training:** Hours or days of uninterrupted GPU/CPU compute cycles.
- **Large-Scale Software Compilation & Rendering:** 3D CAD modeling, game engine builds, and data analytics pipelines.
- **Automated Laboratory Scripts & Server Daemons:** Continuous background data processing.

If a conventional operating system unexpectedly triggers power-saving sleep, screensavers, or display standby during these tasks, training checkpoints are corrupted, socket connections drop, and hours of research work are lost.

### 🚀 Coffee24h as a Core OS Building Block
Within the broader roadmap of a university-centric OS, Coffee24h serves as the **Native Power Orchestrator & Compute Persistence Subsystem**:

```
┌────────────────────────────────────────────────────────────────────────┐
│                   Future University OS (GDU OS Vision)                  │
├────────────────────────────────────────────────────────────────────────┤
│  User Space Applications & Research Tools (IDEs, AI Labs, Compilers)   │
├────────────────────────────────────────────────────────────────────────┤
│  🎯 Coffee24h System Daemon & Subsystem Core:                          │
│     ├─ Low-level Execution State Lock (Kernel Sleep Interception)      │
│     ├─ Real-time Hardware Battery Telemetry & Safeguard (<5%)          │
│     ├─ Background Workstation Security Lock (Win+L / PAM Lock)         │
│     ├─ System Tray Shell Integration & Crash Auto-Recovery             │
│     └─ Microsecond-accurate Harmonic Melodic Audio Alerts              │
├────────────────────────────────────────────────────────────────────────┤
│  OS Kernel / Driver Interfaces (Win32, Darwin IOKit, Linux systemd)    │
└────────────────────────────────────────────────────────────────────────┘
```

By proving that a lightweight, native Go engine can interface directly with OS execution registers while hosting an ultra-responsive, GPU-accelerated DirectComposition UI, Coffee24h establishes a battle-tested foundation for future system software developed at GDU.

---

## ⚡ Core Capabilities & Keep-Awake Modes

| Mode / Feature | Behavior & Purpose | Typical Use Case |
| :--- | :--- | :--- |
| **♾️ Vô Hạn (Until Off)** | Blocks all OS standby, sleep timers, display dimming, and screensavers until manual stop. | Heavy interactive coding, hackathons, lab presentations. |
| **💻 Chạy Nền (System-Only)** | Allows the monitor to turn off to save power and protect OLED/IPS panels, while CPU, network downloads, and compute threads stay at 100%. | Overnight model training, large dataset downloads, video rendering. |
| **⏱️ Hẹn Giờ (Timed Session)** | Quick presets (**15m**, **30m**, **1h**, **2h**, **4h**) with an active countdown timer and 3-tone bell melody (`E5 -> G5 -> C6`). | Pomodoro focus sessions, timed classroom exams, scheduled lab runs. |
| **🛡️ Bảo Vệ Pin (Safeguard)** | Automatically releases wake lock when battery drops below 5% to prevent sudden power outage. | Traveling students, un-plugged library study sessions. |
| **🔒 Khóa Máy (Win + L)** | Instant 1-click workstation lock via native API call (`LockWorkStation`). | Stepping away from the desk in public university libraries. |
| **💎 Bo Tròn 4 Góc (Subpixel)** | Pure DirectComposition 32-bit ARGB with 14px anti-aliased rounded corners and no jagged 1-bit GDI stair-stepping. | Modern, aesthetic glassmorphism UI matching Windows 11/10. |

---

## 📥 Installation Guide (Windows, macOS, Linux)

Coffee24h is distributed in pre-compiled binary packages for all major operating systems. Choose the installation method for your platform below:

### 🪟 1. Windows Installation (Windows 10 / 11)

#### Method A: Full Windows Setup Wizard (Recommended)
1. Download **`coffee24h-installer.exe`** from the [GitHub Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Double-click the installer and follow the prompt.
3. It will automatically install Coffee24h, place high-res GDU shortcuts on your Desktop & Start Menu, and register the uninstaller in Windows Settings.

#### Method B: Standalone Portable Executable (Zero Installation)
1. Download **`coffee24h.exe`** from [GitHub Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Save it to any folder of your choice (e.g. `C:\Tools\Coffee24h\` or your Desktop).
3. Double-click **`coffee24h.exe`** to launch immediately. No administrator rights required!

#### Method C: 1-Click Multi-Setup Manager (From Repository)
If you clone or download the repository, run the built-in multi-setup manager:
- Double-click **`setup.cmd`** or run in PowerShell:
  ```powershell
  powershell -ExecutionPolicy Bypass -File .\setup.ps1 -All
  ```
  This interactive script can configure auto-start on boot (into the System Tray), add `coffee24h` to your global `PATH`, enable `Win + R` run dialog execution, and configure PowerShell profile aliases (`coffee`, `c24`).

---

### 🍎 2. macOS Installation (Apple Silicon & Intel)

Coffee24h provides a universal macOS build that runs natively on both **Apple Silicon (M1/M2/M3/M4)** and **Intel x86_64** Macs.

#### Method A: DMG Drag & Drop Installer (Recommended)
1. Download **`Coffee24h-macOS-Universal.dmg`** from [GitHub Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Double-click the `.dmg` file to mount it.
3. Drag **`Coffee24h.app`** into your **`/Applications`** folder.
4. Launch Coffee24h from Launchpad, Spotlight (`Cmd + Space` -> *Coffee24h*), or Finder.

> [!TIP]
> **Handling macOS Gatekeeper (Unsigned App Notice):**  
> If macOS alerts: *"Coffee24h cannot be opened because it is from an unidentified developer"*:
> 1. Right-click (or Control-click) `Coffee24h.app` in Finder and click **Open**. Click **Open** again in the confirmation prompt.
> 2. Alternatively, remove the quarantine attribute via Terminal:
>    ```bash
>    xattr -cr /Applications/Coffee24h.app
>    ```

#### Method B: Terminal Command-Line Integration
To run Coffee24h directly from zsh/bash:
```bash
# Add alias to your ~/.zshrc or ~/.bashrc
echo 'alias coffee24h="/Applications/Coffee24h.app/Contents/MacOS/Coffee24h"' >> ~/.zshrc
source ~/.zshrc

# Run directly from terminal:
coffee24h --awake --time 2h
```

---

### 🐧 3. Linux Installation (Ubuntu, Debian, Fedora, Arch)

#### Method A: Debian / Ubuntu Package (`.deb`)
For Ubuntu, Debian, Linux Mint, Pop!_OS, and other Debian-based distributions:
1. Download **`coffee24h_1.0.0_amd64.deb`** from [GitHub Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Install via APT or dpkg:
   ```bash
   sudo apt install ./coffee24h_1.0.0_amd64.deb
   # Or using dpkg:
   sudo dpkg -i coffee24h_1.0.0_amd64.deb
   ```
3. Coffee24h will be installed into `/usr/bin/coffee24h` with a desktop entry and custom icon registered in your application launcher.

#### Method B: Standalone Portable AppImage (Universal Linux)
Runs on almost any modern 64-bit Linux distribution (Ubuntu, Fedora, Arch, openSUSE, Manjaro):
1. Download **`Coffee24h-x86_64.AppImage`** from [GitHub Releases](https://github.com/nynwadhyhcz8/Coffee24h/releases).
2. Grant execution permissions:
   ```bash
   chmod +x Coffee24h-x86_64.AppImage
   ```
3. Run the AppImage:
   ```bash
   ./Coffee24h-x86_64.AppImage
   ```

#### Method C: Manual Binary Installation (Arch / Fedora / Generic)
If you build or download the standalone `coffee24h` Linux executable:
```bash
# Copy binary to system PATH
sudo install -m 755 coffee24h /usr/local/bin/coffee24h

# (Optional) Install desktop entry and icon for system menu
sudo cp build/linux/coffee24h.desktop /usr/share/applications/
sudo cp build/linux/coffee24h.png /usr/share/icons/hicolor/256x256/apps/
```

---

### 📦 4. NPM & Node.js (CLI & Web Preview)

For web developers, Node.js engineers, or users who want to run Coffee24h via `npm` or test the interface in a browser without compiling Go binaries:

#### Option A: Run via NPM CLI Script
Clone the repository and run immediately via Node.js:
```bash
git clone https://github.com/nynwadhyhcz8/Coffee24h.git
cd Coffee24h

# Run Coffee24h directly (auto-detects compiled binary or launches Web UI):
npm start

# Pass CLI flags through npm:
npm start -- --awake --time 2h
npm start -- --student
```

#### Option B: Global NPM CLI Installation
Install the launcher shim globally across your entire system:
```bash
# From within the cloned folder:
npm install -g .

# Now execute 'coffee24h' or 'coffee-gdu' from any directory:
coffee24h --awake
coffee24h --tray
```

#### Option C: Live Browser UI Preview (Zero Go Dependencies)
Perfect for testing the GDU Glassmorphism interface, interactive timer countdown, and Web Audio API synthesizer directly in Chrome, Edge, Safari, or Firefox:
```bash
# Start Vite development server from root:
npm run dev

# Or build production static web assets:
npm run build
npm run preview
```
Visit `http://localhost:5173` to test the full frontend experience.

---

## 🏗️ Low-Level OS Architecture

Coffee24h utilizes a **Hybrid Systems Architecture** that pairs the raw speed and direct Win32/POSIX syscall access of **Go** with the visual elegance of **Vite + CSS3 Glassmorphism**:

```
┌────────────────────────────────────────────────────────────┐
│                    Frontend UI Layer                       │
│    Vite + Vanilla JS + CSS3 Glassmorphism (350x350 1:1)    │
│       - Dual Edition Switcher (Lecturer ⇄ Student)         │
│       - Web Audio Harmonic Synthesizer (Offline Chimes)    │
│       - Bilingual Reactive i18n (Tiếng Việt & English)     │
└────────────────────────────▲───────────────────────────────┘
                             │ Wails v2 Bidirectional IPC Bridge
┌────────────────────────────▼───────────────────────────────┐
│                   Go Core Subsystem Layer                  │
│    - The Shield: Thread-safe atomic Mutex (sync.Mutex)     │
│    - Pinned OS Thread Execution (runtime.LockOSThread)     │
│    - Deterministic State Machine (AwakeState Engine)       │
├────────────────────────────────────────────────────────────┤
│                Native OS Kernel Interfacing                │
│    Windows (Win32):                                        │
│      • SetThreadExecutionState (ES_CONTINUOUS | ...)       │
│      • Shell_NotifyIconW (System Tray with TaskbarCreated) │
│      • LockWorkStation (Instant physical security)         │
│      • GetSystemPowerStatus (Battery telemetry)            │
│    macOS (Darwin):                                         │
│      • caffeinate -d -i -m -u & pmset battery parsing      │
│    Linux (X11 / Wayland):                                  │
│      • systemd-inhibit / org.freedesktop.ScreenSaver       │
└────────────────────────────────────────────────────────────┘
```

---

## 🖥️ Cross-Platform Hardware & Power Mapping

| Operating System | Sleep Interception Mechanism | Display State Handling | Battery Telemetry Source | Workstation Security Lock |
| :--- | :--- | :--- | :--- | :--- |
| **🪟 Windows (10 / 11)** | `SetThreadExecutionState` with pinned OS thread | `ES_DISPLAY_REQUIRED` flag toggling | Win32 `GetSystemPowerStatus` API | Win32 `LockWorkStation` |
| **🍎 macOS (Apple Silicon / Intel)** | `caffeinate -i -d` subprocess orchestration | Dynamic flag switching (`-d` / `-i`) | `pmset -g batt` POSIX pipe | `pmset displaysleepnow` |
| **🐧 Linux (X11 / Wayland)** | `systemd-inhibit` / D-Bus Screensaver API | D-Bus lock acquisition | Linux kernel `/sys/class/power_supply/` | `loginctl lock-session` |

---

## 🛠️ Build from Source & Packaging

### Prerequisites for Developers
- **Go** (1.20+): [golang.org/dl](https://golang.org/dl/)
- **Node.js** (npm LTS): [nodejs.org](https://nodejs.org/)
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Platform Dependencies**:
  - *Windows:* WebView2 Runtime (Pre-installed on Win 10/11)
  - *macOS:* Xcode Command Line Tools (`xcode-select --install`)
  - *Linux:* `sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev` (or `webkit2gtk-4.1-dev`)

### 🚀 Building on Windows
```powershell
cd coffeeGDU
# Option 1: Standard Wails build
wails build

# Option 2: Build NSIS Setup Installer
wails build -nsis
```

### 🍎 Building & Packaging on macOS
```bash
cd coffeeGDU
# Build universal binary & package into .dmg
chmod +x ./build/darwin/build-dmg.sh
./build/darwin/build-dmg.sh
```

### 🐧 Building & Packaging on Linux
```bash
cd coffeeGDU
# Build Debian/Ubuntu .deb package:
chmod +x ./build/linux/build-deb.sh
./build/linux/build-deb.sh 1.0.0

# Build universal .AppImage:
chmod +x ./build/linux/build-appimage.sh
./build/linux/build-appimage.sh
```

---

## ⚙️ 1-Click Multi-Setup Manager

To deploy Coffee24h effortlessly on any Windows workstation, double-click **`setup.cmd`** (or execute `setup.ps1` in PowerShell):

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

Coffee24h includes a full-featured terminal CLI, allowing lab administrators, students, and automation scripts to control wake states without opening the GUI:

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
```

### CLI Quick Examples:
```bash
coffee24h --awake                  # Launch and immediately keep laptop awake
coffee24h --tray                   # Start silently into system tray
coffee24h --time 2h --autolock     # Keep awake for 2 hours, then lock machine
coffee24h -s --time 1h             # 1 hour background compute (screen can sleep)
coffee24h --student                # Launch directly in GDU Student Edition
```

---

## 👨‍💻 Author & Acknowledgments

- **Creator & Lead Engineer:** **Neith Do** ([@nynwadhyhcz8](https://github.com/nynwadhyhcz8) / `16av112thien@gmail.com`)
- **Institution:** Dedicated with passion to **Gia Dinh University (GDU)** — Faculty, Researchers, and Students.
- **Architectural Reference:** Detailed troubleshooting and code patterns are documented in [`DEVELOPER_GUIDE.md`](./DEVELOPER_GUIDE.md).

<p align="center">
  <strong>Made with ❤️, ☕, and ambition for Gia Dinh University (GDU)</strong><br />
  <em>"Chọn Gia Định — Chọn Thành Công"</em>
</p>
