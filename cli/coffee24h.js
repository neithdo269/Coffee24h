#!/usr/bin/env node

/**
 * Coffee24h — Cross-Platform CLI Launcher & NPM Runner
 * Part of the Gia Dinh University (GDU) Keep-Awake Project
 * Author: Neith Do (@nynwadhyhcz8)
 */

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const platform = os.platform();
const args = process.argv.slice(2);
const rootDir = path.resolve(__dirname, '..');

// Handle --help and --version directly in terminal
if (args.includes('-h') || args.includes('--help') || args.includes('/?') || args.includes('help')) {
  console.log(`
☕ Coffee24h — Keep-Awake Utility & System Power Daemon (GDU Edition)

Usage:
  coffee24h [flags]
  npm start -- [flags]

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
  npm start                       # Launch GUI normally
  npm start -- --tray             # Start silently into system tray
  npm start -- --awake --time 2h  # Keep awake for 2 hours with timer
  npm start -- --student          # Launch directly with GDU Student Edition
`);
  process.exit(0);
}

if (args.includes('-v') || args.includes('--version') || args.includes('version')) {
  console.log('Coffee24h (GDU Edition) v1.0.0');
  process.exit(0);
}

// Locate compiled native binary
let binPath = null;

if (platform === 'win32') {
  const candidates = [
    path.join(rootDir, 'coffeeGDU', 'coffee24h.exe'),
    path.join(rootDir, 'coffeeGDU', 'build', 'bin', 'coffee24h.exe'),
    path.join(rootDir, 'coffee24h.exe')
  ];
  binPath = candidates.find(p => fs.existsSync(p));
} else if (platform === 'darwin') {
  const candidates = [
    path.join(rootDir, 'coffeeGDU', 'build', 'bin', 'Coffee24h.app', 'Contents', 'MacOS', 'Coffee24h'),
    path.join(rootDir, 'coffeeGDU', 'build', 'bin', 'coffee24h'),
    path.join(rootDir, 'coffee24h')
  ];
  binPath = candidates.find(p => fs.existsSync(p));
} else {
  // Linux / Unix
  const candidates = [
    path.join(rootDir, 'coffeeGDU', 'build', 'bin', 'coffee24h'),
    path.join(rootDir, 'coffeeGDU', 'coffee24h'),
    path.join(rootDir, 'coffee24h')
  ];
  binPath = candidates.find(p => fs.existsSync(p));
}

if (binPath && fs.existsSync(binPath)) {
  const child = spawn(binPath, args, {
    cwd: path.dirname(binPath),
    stdio: 'inherit',
    windowsHide: false
  });
  child.on('exit', (code) => process.exit(code || 0));
} else {
  // Fallback to launching the frontend Web Preview via Vite
  console.log('☕ [Coffee24h] Compiled binary not found. Launching Browser UI Preview via Vite...');
  const frontendDir = path.join(rootDir, 'coffeeGDU', 'frontend');
  const npmCmd = platform === 'win32' ? 'npm.cmd' : 'npm';
  
  const dev = spawn(npmCmd, ['run', 'dev'], {
    cwd: frontendDir,
    stdio: 'inherit',
    shell: true
  });
  dev.on('exit', (code) => process.exit(code || 0));
}
