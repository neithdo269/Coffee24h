package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// CLIOptions encapsulates command-line arguments parsed at startup
type CLIOptions struct {
	ShowHelp    bool
	ShowVersion bool
	StartTray   bool
	AwakeNow    bool
	TimedMin    int
	SystemOnly  bool
	NoSafeguard bool
	AutoLock    bool
	EcoCup      bool
}

func parseCLI() CLIOptions {
	var opts CLIOptions
	for i := 1; i < len(os.Args); i++ {
		arg := strings.ToLower(os.Args[i])
		switch {
		case arg == "-h" || arg == "--help" || arg == "/?" || arg == "help":
			opts.ShowHelp = true
		case arg == "-v" || arg == "--version" || arg == "version":
			opts.ShowVersion = true
		case arg == "-m" || arg == "--tray" || arg == "--minimized" || arg == "--min" || arg == "-min":
			opts.StartTray = true
		case arg == "-a" || arg == "--awake" || arg == "--on" || arg == "on" || arg == "awake":
			opts.AwakeNow = true
		case arg == "-s" || arg == "--system-only" || arg == "--system":
			opts.SystemOnly = true
		case arg == "--no-safeguard":
			opts.NoSafeguard = true
		case arg == "--autolock" || arg == "--lock" || arg == "-l":
			opts.AutoLock = true
		case arg == "-e" || arg == "--student" || arg == "--recycle" || arg == "--leaf" || arg == "--eco":
			opts.EcoCup = true
		case arg == "--lecturer" || arg == "--teacher" || arg == "--gdu":
			opts.EcoCup = false
		case arg == "-t" || arg == "--time" || arg == "--duration":
			if i+1 < len(os.Args) {
				opts.TimedMin = parseDuration(os.Args[i+1])
				i++
			}
		case strings.HasPrefix(arg, "--time="):
			opts.TimedMin = parseDuration(strings.TrimPrefix(arg, "--time="))
		case strings.HasPrefix(arg, "-t="):
			opts.TimedMin = parseDuration(strings.TrimPrefix(arg, "-t="))
		}
	}
	return opts
}

func parseDuration(s string) int {
	s = strings.TrimSpace(strings.ToLower(s))
	if strings.HasSuffix(s, "h") {
		v, _ := strconv.Atoi(strings.TrimSuffix(s, "h"))
		return v * 60
	}
	if strings.HasSuffix(s, "m") {
		v, _ := strconv.Atoi(strings.TrimSuffix(s, "m"))
		return v
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return 0
}

func attachConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole := kernel32.NewProc("AttachConsole")
	const ATTACH_PARENT_PROCESS = ^uintptr(0) // (DWORD)-1
	r, _, _ := procAttachConsole.Call(ATTACH_PARENT_PROCESS)
	if r != 0 {
		conout, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
		if err == nil {
			os.Stdout = conout
			os.Stderr = conout
		}
	}
}

func printHelp() {
	attachConsole()
	msg := `
  ☕ Coffee24h - Keep-Awake Utility for Gia Dinh University (GDU)
  Version: 1.0.0 | Author: Neith Do

  USAGE:
    coffee24h [flags]

  FLAGS:
    -h, --help            Show this help reference
    -v, --version         Show application version
    -m, --tray            Start directly minimized into System Tray ("Show hidden icons")
    -a, --awake           Immediately activate Keep-Awake mode upon launch
    -e, --student         Trigger GDU Student Edition (Recycle Leaf Cup)
    --lecturer            Trigger GDU Lecturer Edition (Blue Cup with GDU Logo)
    -t, --time <dur>      Start a timed session (e.g. 15m, 30m, 1h, 2h, 4h)
    -s, --system-only     Keep CPU/downloads awake with display sleep enabled
    -l, --autolock        Auto-lock workstation (Win + L) when session expires
    --no-safeguard        Disable battery safeguard (<5% auto-stop)

  EXAMPLES:
    coffee24h                    # Launch GUI normally (Lecturer Blue Cup with GDU Logo)
    coffee24h --student          # Launch in GDU Student Edition (Recycle Leaf Cup)
    coffee24h --tray             # Start silently in "Show hidden icons"
    coffee24h --awake            # Launch and immediately keep laptop awake
    coffee24h --time 2h          # Keep awake for 2 hours with timer
    coffee24h -e --time 1h       # Student session for 1 hour
`
	fmt.Println(msg)
}

func printVersion() {
	attachConsole()
	fmt.Println("Coffee24h v1.0.0 (GDU Edition) - Native Win32 Keep-Awake Utility")
}
