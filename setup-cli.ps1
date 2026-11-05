<#
.SYNOPSIS
    Coffee24h Global Command-Line Setup Script
.DESCRIPTION
    Integrates 'coffee24h' into your Windows terminal environment permanently.
    Allows running 'coffee24h' from any directory in PowerShell, CMD, or Terminal.
#>

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinDir = Join-Path $ScriptDir "coffeeGDU"
$ExePath = Join-Path $BinDir "coffee24h.exe"

if (-not (Test-Path $ExePath)) {
    Write-Host "[-] coffee24h.exe not found at $ExePath" -ForegroundColor Red
    exit 1
}

Write-Host "☕ Setting up Coffee24h Command Line Interface..." -ForegroundColor Cyan

# 1. Add coffeeGDU to User PATH permanently
$userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$BinDir*") {
    $newUserPath = "$BinDir;$userPath"
    [System.Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    Write-Host "[+] Added $BinDir to User PATH permanently." -ForegroundColor Green
} else {
    Write-Host "[i] $BinDir is already in User PATH." -ForegroundColor DarkGray
}

# 2. Deploy instant launcher shim to WinGet Links (already in active system PATH)
$winGetLinks = "$env:LOCALAPPDATA\Microsoft\WinGet\Links"
if (Test-Path $winGetLinks) {
    $shimCmd = Join-Path $winGetLinks "coffee24h.cmd"
    $shimContent = "@echo off`r`n`"$ExePath`" %*"
    [System.IO.File]::WriteAllText($shimCmd, $shimContent, [System.Text.Encoding]::ASCII)
    Write-Host "[+] Instant launcher shim deployed to $shimCmd" -ForegroundColor Green
}

# 3. Update current session PATH so it works immediately
if ($env:Path -notlike "*$BinDir*") {
    $env:Path = "$BinDir;$env:Path"
}

Write-Host ""
Write-Host "==========================================================" -ForegroundColor Green
Write-Host " ✨ Coffee24h CLI Setup Complete! " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
Write-Host "You can now run 'coffee24h' from anywhere in terminal!" -ForegroundColor Cyan
Write-Host ""
Write-Host "Commands you can try:" -ForegroundColor Yellow
Write-Host "  coffee24h                Launch GUI" -ForegroundColor White
Write-Host "  coffee24h --tray         Start silently in System Tray" -ForegroundColor White
Write-Host "  coffee24h --awake        Start keep-awake mode immediately" -ForegroundColor White
Write-Host "  coffee24h --time 2h      Start 2-hour keep-awake timer" -ForegroundColor White
Write-Host "  coffee24h --system-only  Keep CPU alive, allow screen to sleep" -ForegroundColor White
Write-Host "  coffee24h --help         Show full CLI flag reference" -ForegroundColor White
Write-Host ""
