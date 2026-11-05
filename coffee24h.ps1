# PowerShell CLI Launcher for Coffee24h
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
& "$ScriptDir\coffeeGDU\coffee24h.exe" @args
