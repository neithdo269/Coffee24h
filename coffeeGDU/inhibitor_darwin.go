//go:build darwin

package main

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type DarwinInhibitor struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	isActive bool
}

func newPlatformInhibitor() PlatformInhibitor {
	return &DarwinInhibitor{}
}

func (d *DarwinInhibitor) Inhibit(keepDisplay bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd != nil && d.cmd.Process != nil {
		_ = d.cmd.Process.Kill()
	}

	args := []string{"-i"}
	if keepDisplay {
		args = append(args, "-d")
	}

	cmd := exec.Command("caffeinate", args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	d.cmd = cmd
	d.isActive = true
	return nil
}

func (d *DarwinInhibitor) Release() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd != nil && d.cmd.Process != nil {
		_ = d.cmd.Process.Kill()
		d.cmd = nil
	}
	d.isActive = false
	return nil
}

func (d *DarwinInhibitor) IsActive() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isActive
}

// PlaySystemSound plays an alert sound on macOS using native system audio
func PlaySystemSound(soundType uint) {
	// 0x10 is critical stop/alert, 0x40 is info/chime
	soundFile := "/System/Library/Sounds/Ping.aiff"
	if soundType == 0x10 {
		soundFile = "/System/Library/Sounds/Basso.aiff"
	}
	_ = exec.Command("afplay", soundFile).Start()
}

// LockScreen locks the macOS session, with fallback to display sleep
func LockScreen() error {
	if err := exec.Command("/System/Library/CoreServices/Menu Extras/User.menu/Contents/Resources/CGSession", "-suspend").Run(); err == nil {
		return nil
	}
	return exec.Command("pmset", "displaysleepnow").Run()
}

func getPlatformBatteryInfo() BatteryInfo {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return BatteryInfo{HasBattery: false, PowerSource: "AC / Unknown"}
	}
	text := string(out)
	hasBattery := strings.Contains(text, "InternalBattery")
	pct := 100
	isCharging := strings.Contains(text, "charging")
	powerSource := "AC Adapter"
	if strings.Contains(text, "Battery Power") {
		powerSource = "Battery"
	}

	// Simple percent extraction
	if idx := strings.Index(text, "%"); idx != -1 {
		start := idx - 1
		for start >= 0 && text[start] >= '0' && text[start] <= '9' {
			start--
		}
		if p, err := strconv.Atoi(text[start+1 : idx]); err == nil {
			pct = p
		}
	}

	return BatteryInfo{
		HasBattery:  hasBattery,
		Percentage:  pct,
		IsCharging:  isCharging,
		PowerSource: powerSource,
	}
}
