//go:build linux

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type LinuxInhibitor struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	isActive bool
}

func newPlatformInhibitor() PlatformInhibitor {
	return &LinuxInhibitor{}
}

func (l *LinuxInhibitor) Inhibit(keepDisplay bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cmd != nil && l.cmd.Process != nil {
		_ = l.cmd.Process.Kill()
	}

	what := "idle:sleep"
	if keepDisplay {
		what = "idle:sleep:handle-lid-switch"
	}

	cmd := exec.Command("systemd-inhibit", "--what="+what, "--who=Coffee24h", "--why=User requested stay awake", "sleep", "infinity")
	if err := cmd.Start(); err != nil {
		fallback := exec.Command("xdg-screensaver", "suspend", "0x0")
		_ = fallback.Start()
		l.cmd = fallback
	} else {
		l.cmd = cmd
	}

	l.isActive = true
	return nil
}

func (l *LinuxInhibitor) Release() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cmd != nil && l.cmd.Process != nil {
		_ = l.cmd.Process.Kill()
		l.cmd = nil
	}
	l.isActive = false
	return nil
}

func (l *LinuxInhibitor) IsActive() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.isActive
}

// PlaySystemSound plays an alert sound on Linux using freedesktop sound themes or terminal bell
func PlaySystemSound(soundType uint) {
	soundName := "complete"
	if soundType == 0x10 {
		soundName = "dialog-warning"
	}
	cmd := exec.Command("canberra-gtk-play", "-i", soundName)
	if err := cmd.Start(); err != nil {
		_ = exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/"+soundName+".oga").Start()
	}
}

// LockScreen locks the Linux desktop session via loginctl (systemd) or xdg-screensaver
func LockScreen() error {
	if err := exec.Command("loginctl", "lock-session").Run(); err == nil {
		return nil
	}
	return exec.Command("xdg-screensaver", "lock").Run()
}

func getPlatformBatteryInfo() BatteryInfo {
	// Check /sys/class/power_supply/BAT*
	matches, _ := filepath.Glob("/sys/class/power_supply/BAT*")
	if len(matches) == 0 {
		return BatteryInfo{HasBattery: false, PowerSource: "AC / Unknown"}
	}

	batDir := matches[0]
	capacityBytes, _ := os.ReadFile(filepath.Join(batDir, "capacity"))
	statusBytes, _ := os.ReadFile(filepath.Join(batDir, "status"))

	pct, _ := strconv.Atoi(strings.TrimSpace(string(capacityBytes)))
	status := strings.TrimSpace(string(statusBytes))
	isCharging := status == "Charging"
	powerSource := "Battery"
	if isCharging || status == "Full" {
		powerSource = "AC Adapter"
	}

	return BatteryInfo{
		HasBattery:  true,
		Percentage:  pct,
		IsCharging:  isCharging,
		PowerSource: powerSource,
	}
}
