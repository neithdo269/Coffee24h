<#
.SYNOPSIS
    Coffee24h Multi-Setup Manager
.DESCRIPTION
    Provides various ways to set up and integrate Coffee24h on Windows:
    - Global CLI in Terminal (PATH + WinGet Links)
    - Auto-Start on Windows Login (Silent to System Tray)
    - Desktop Shortcut with Custom Icon
    - Start Menu Shortcut (Searchable in Windows Search)
    - Run Dialog integration (Win + R -> 'coffee24h')
    - PowerShell Profile Aliases ('coffee', 'c24')
    - Complete Uninstaller
#>

param(
    [switch]$All,
    [switch]$Cli,
    [switch]$Startup,
    [switch]$DesktopShortcut,
    [switch]$StartMenu,
    [switch]$RunDialog,
    [switch]$PowerShellProfile,
    [switch]$Uninstall,
    [switch]$Silent
)

$ErrorActionPreference = "Continue"

# Resolve directories
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
if (Test-Path (Join-Path $ScriptDir "coffeeGDU\coffee24h.exe")) {
    $BinDir = Join-Path $ScriptDir "coffeeGDU"
} elseif (Test-Path (Join-Path $ScriptDir "coffee24h.exe")) {
    $BinDir = $ScriptDir
} else {
    Write-Host "[-] Could not locate coffee24h.exe" -ForegroundColor Red
    exit 1
}

$ExePath = Join-Path $BinDir "coffee24h.exe"
$IconPath = Join-Path $BinDir "CoffeeGDU icons\Coffee-24h-GDU-64Gperpackage-icon.ico"
if (-not (Test-Path $IconPath)) {
    $IconPath = Join-Path $BinDir "build\windows\icon.ico"
}

function Show-Header {
    Write-Host "==============================================================" -ForegroundColor Cyan
    Write-Host "   Coffee24h - Multi-Setup Manager (Various Setup Options)" -ForegroundColor White
    Write-Host "  Keep-Awake Utility for Gia Dinh University (GDU)" -ForegroundColor DarkGray
    Write-Host "==============================================================" -ForegroundColor Cyan
    Write-Host "Target Executable: $ExePath" -ForegroundColor DarkGray
    Write-Host ""
}

# --- Action 1: CLI (PATH + WinGet Link) ---
function Install-CliSetup {
    Write-Host "[*] 1. Setting up Command-Line Interface (CLI)..." -ForegroundColor Yellow

    # Add to User PATH
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$BinDir*") {
        $newUserPath = "$BinDir;$userPath"
        [System.Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
        Write-Host "  [+] Added $BinDir to User PATH." -ForegroundColor Green
    } else {
        Write-Host "  [i] $BinDir is already in User PATH." -ForegroundColor DarkGray
    }

    # Deploy WinGet Links shim (already in active system PATH)
    $winGetLinks = "$env:LOCALAPPDATA\Microsoft\WinGet\Links"
    if (Test-Path $winGetLinks) {
        $shimCmd = Join-Path $winGetLinks "coffee24h.cmd"
        $shimContent = "@echo off`r`n`"$ExePath`" %*"
        [System.IO.File]::WriteAllText($shimCmd, $shimContent, [System.Text.Encoding]::ASCII)
        Write-Host "  [+] Launcher shim created in WinGet\Links (instant terminal access)." -ForegroundColor Green
    }

    # Update active session
    if ($env:Path -notlike "*$BinDir*") {
        $env:Path = "$BinDir;$env:Path"
    }
}

# --- Action 2: Windows Auto-Start on Boot ---
function Install-Startup {
    Write-Host "[*] 2. Configuring Auto-Start on Windows Boot..." -ForegroundColor Yellow
    $runKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
    $runCmd = "`"$ExePath`" --tray"
    Set-ItemProperty -Path $runKey -Name "Coffee24h" -Value $runCmd
    Write-Host "  [+] Registered in Windows Startup: Starts silently in System Tray on login." -ForegroundColor Green
}

function Remove-Startup {
    $runKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
    Remove-ItemProperty -Path $runKey -Name "Coffee24h" -ErrorAction SilentlyContinue
    Write-Host "  [-] Removed Coffee24h from Windows Startup." -ForegroundColor DarkYellow
}

# --- Action 3: Desktop Shortcut ---
function Install-DesktopShortcut {
    Write-Host "[*] 3. Creating Desktop Shortcut..." -ForegroundColor Yellow
    $desktopDir = [Environment]::GetFolderPath('Desktop')
    $shortcutPath = Join-Path $desktopDir "Coffee24h.lnk"

    $wsh = New-Object -ComObject WScript.Shell
    $shortcut = $wsh.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $ExePath
    if (Test-Path $IconPath) {
        $shortcut.IconLocation = "$IconPath, 0"
    }
    $shortcut.WorkingDirectory = $BinDir
    $shortcut.Description = "Coffee24h - Keep-Awake Utility for GDU"
    $shortcut.Save()
    Write-Host "  [+] Desktop shortcut created at $shortcutPath" -ForegroundColor Green
}

# --- Action 4: Start Menu Shortcut ---
function Install-StartMenu {
    Write-Host "[*] 4. Creating Start Menu Shortcut..." -ForegroundColor Yellow
    $startMenuPrograms = Join-Path ([Environment]::GetFolderPath('StartMenu')) "Programs"
    $shortcutPath = Join-Path $startMenuPrograms "Coffee24h.lnk"

    $wsh = New-Object -ComObject WScript.Shell
    $shortcut = $wsh.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $ExePath
    if (Test-Path $IconPath) {
        $shortcut.IconLocation = "$IconPath, 0"
    }
    $shortcut.WorkingDirectory = $BinDir
    $shortcut.Description = "Coffee24h - Keep-Awake Utility for GDU"
    $shortcut.Save()
    Write-Host "  [+] Start Menu shortcut created: searchable as 'Coffee24h' in Windows search." -ForegroundColor Green
}

# --- Action 5: Run Dialog (Win + R) ---
function Install-RunDialog {
    Write-Host "[*] 5. Registering Windows Run Dialog (Win + R)..." -ForegroundColor Yellow
    $appPathsKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\App Paths\coffee24h.exe"
    if (-not (Test-Path $appPathsKey)) {
        New-Item -Path $appPathsKey -Force | Out-Null
    }
    Set-ItemProperty -Path $appPathsKey -Name "(Default)" -Value $ExePath
    Set-ItemProperty -Path $appPathsKey -Name "Path" -Value $BinDir
    Write-Host "  [+] App Paths registered: Press Win + R, type 'coffee24h', and hit Enter!" -ForegroundColor Green
}

# --- Action 6: PowerShell Profile Aliases ---
function Install-PowerShellProfile {
    Write-Host "[*] 6. Configuring PowerShell Profile Shortcuts..." -ForegroundColor Yellow
    $profileDir = Split-Path -Parent $PROFILE
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $PROFILE)) {
        New-Item -ItemType File -Path $PROFILE -Force | Out-Null
    }

    $aliasCode = "`r`n# >>> Coffee24h Keep-Awake Shortcuts >>>`r`nfunction coffee { & `"$ExePath`" @args }`r`nfunction c24 { & `"$ExePath`" @args }`r`n# <<< Coffee24h Keep-Awake Shortcuts <<<`r`n"

    Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force -ErrorAction SilentlyContinue
    $profileText = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
    if ($profileText -notlike "*Coffee24h Keep-Awake Shortcuts*") {
        Add-Content -Path $PROFILE -Value $aliasCode
        Write-Host "  [+] Added 'coffee' and 'c24' command aliases to $PROFILE" -ForegroundColor Green
    } else {
        Write-Host "  [i] Aliases already present in $PROFILE" -ForegroundColor DarkGray
    }
}

# --- Action 7: Full Uninstall ---
function Uninstall-All {
    Write-Host "[*] Removing all Coffee24h shortcuts and integrations..." -ForegroundColor Yellow

    # Remove Startup
    Remove-Startup

    # Remove Desktop Shortcut
    $desktopShortcut = Join-Path ([Environment]::GetFolderPath('Desktop')) "Coffee24h.lnk"
    if (Test-Path $desktopShortcut) {
        Remove-Item $desktopShortcut -Force
        Write-Host "  [-] Removed Desktop shortcut." -ForegroundColor DarkYellow
    }

    # Remove Start Menu Shortcut
    $startMenuShortcut = Join-Path ([Environment]::GetFolderPath('StartMenu')) "Programs\Coffee24h.lnk"
    if (Test-Path $startMenuShortcut) {
        Remove-Item $startMenuShortcut -Force
        Write-Host "  [-] Removed Start Menu shortcut." -ForegroundColor DarkYellow
    }

    # Remove Run Dialog key
    $appPathsKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\App Paths\coffee24h.exe"
    if (Test-Path $appPathsKey) {
        Remove-Item -Path $appPathsKey -Recurse -Force
        Write-Host "  [-] Removed Run dialog integration." -ForegroundColor DarkYellow
    }

    # Remove WinGet link shim
    $shimCmd = "$env:LOCALAPPDATA\Microsoft\WinGet\Links\coffee24h.cmd"
    if (Test-Path $shimCmd) {
        Remove-Item $shimCmd -Force
        Write-Host "  [-] Removed WinGet Links launcher shim." -ForegroundColor DarkYellow
    }

    Write-Host ""
    Write-Host "[+] All integrations removed successfully." -ForegroundColor Green
}

# Handle command-line arguments
if ($All) {
    Show-Header
    Install-CliSetup
    Install-Startup
    Install-DesktopShortcut
    Install-StartMenu
    Install-RunDialog
    Install-PowerShellProfile
    Write-Host ""
    Write-Host " All 6 setups completed successfully!" -ForegroundColor Green
    exit 0
}

if ($Uninstall) {
    Show-Header
    Uninstall-All
    exit 0
}

if ($Cli -or $Startup -or $DesktopShortcut -or $StartMenu -or $RunDialog -or $PowerShellProfile) {
    Show-Header
    if ($Cli) { Install-CliSetup }
    if ($Startup) { Install-Startup }
    if ($DesktopShortcut) { Install-DesktopShortcut }
    if ($StartMenu) { Install-StartMenu }
    if ($RunDialog) { Install-RunDialog }
    if ($PowerShellProfile) { Install-PowerShellProfile }
    Write-Host ""
    Write-Host " Selected setup actions completed!" -ForegroundColor Green
    exit 0
}

# Interactive Menu Mode
Show-Header
Write-Host "Select an option to configure:" -ForegroundColor White
Write-Host ""
Write-Host "  [1]  Full Setup (All Recommended Setups in 1-Click)" -ForegroundColor Green
Write-Host "  [2]  Terminal CLI Setup (Add to PATH & WinGet Links)" -ForegroundColor White
Write-Host "  [3]  Auto-Start on Windows Login (Silent to System Tray)" -ForegroundColor White
Write-Host "  [4]  Create Desktop Shortcut (Custom GDU Icon)" -ForegroundColor White
Write-Host "  [5]  Create Start Menu Shortcut (Windows Searchable)" -ForegroundColor White
Write-Host "  [6]  Win + R Run Dialog ('coffee24h')" -ForegroundColor White
Write-Host "  [7]  PowerShell Profile Aliases ('coffee', 'c24')" -ForegroundColor White
Write-Host "  [8]  Remove / Uninstall All Integrations" -ForegroundColor DarkYellow
Write-Host "  [Q]  Quit" -ForegroundColor DarkGray
Write-Host ""

$choice = Read-Host "Enter your choice (1-8 or Q)"
Write-Host ""

switch ($choice) {
    "1" {
        Install-CliSetup
        Install-Startup
        Install-DesktopShortcut
        Install-StartMenu
        Install-RunDialog
        Install-PowerShellProfile
        Write-Host "`n All 6 setups completed successfully!" -ForegroundColor Green
    }
    "2" { Install-CliSetup }
    "3" { Install-Startup }
    "4" { Install-DesktopShortcut }
    "5" { Install-StartMenu }
    "6" { Install-RunDialog }
    "7" { Install-PowerShellProfile }
    "8" { Uninstall-All }
    "Q" { Write-Host "Goodbye!" -ForegroundColor DarkGray; exit 0 }
    default { Write-Host "Invalid choice." -ForegroundColor Red }
}
