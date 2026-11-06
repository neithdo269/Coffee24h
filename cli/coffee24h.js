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
    stdio: 'inherit',
    windowsHide: false
  });
  child.on('exit', (code) => process.exit(code || 0));
} else {
  // Fallback to launching the frontend Web Preview via Vite
  console.log('☕ [Coffee24h] Launching Browser UI Preview via Vite...');
  const frontendDir = path.join(rootDir, 'coffeeGDU', 'frontend');
  const npmCmd = platform === 'win32' ? 'npm.cmd' : 'npm';
  
  const dev = spawn(npmCmd, ['run', 'dev'], {
    cwd: frontendDir,
    stdio: 'inherit',
    shell: true
  });
  dev.on('exit', (code) => process.exit(code || 0));
}
