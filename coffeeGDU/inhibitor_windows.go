//go:build windows

/*
 * Coffee24h - Part of the GDU (Gia Định University) Window Manager Project & OS Accessories.
 * Developed & Maintained by Neith Do.
 * Released under Free & Open Source principles (GDU Open License).
 */

package main

import (
	"syscall"
	"unsafe"
)

// Windows Execution State Flags
const (
	ES_SYSTEM_REQUIRED  = 0x00000001
	ES_DISPLAY_REQUIRED = 0x00000002
	ES_CONTINUOUS       = 0x80000000
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecState   = kernel32.NewProc("SetThreadExecutionState")
	procGetSystemPowerStatus = kernel32.NewProc("GetSystemPowerStatus")

	user32              = syscall.NewLazyDLL("user32.dll")
	procMessageBeep     = user32.NewProc("MessageBeep")
	procLockWorkStation = user32.NewProc("LockWorkStation")
)

type SYSTEM_POWER_STATUS struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	SystemStatusFlag    byte
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

type WindowsInhibitor struct {
	isActive bool
}

func newPlatformInhibitor() PlatformInhibitor {
	return &WindowsInhibitor{}
}

func (w *WindowsInhibitor) Inhibit(keepDisplay bool) error {
	var flags uintptr = ES_CONTINUOUS | ES_SYSTEM_REQUIRED
	if keepDisplay {
		flags |= ES_DISPLAY_REQUIRED
	}

	r1, _, err := procSetThreadExecState.Call(flags)
	if r1 == 0 {
		return err
	}
	w.isActive = true
	return nil
}

func (w *WindowsInhibitor) Release() error {
	r1, _, err := procSetThreadExecState.Call(ES_CONTINUOUS)
	if r1 == 0 {
		return err
	}
	w.isActive = false
	return nil
}

func (w *WindowsInhibitor) IsActive() bool {
	return w.isActive
}

// PlaySystemSound plays a native Windows sound (0x00 = Beep, 0x10 = Critical Stop, 0x40 = Info/Chime)
func PlaySystemSound(soundType uint) {
	_, _, _ = procMessageBeep.Call(uintptr(soundType))
}

func LockScreen() error {
	r1, _, err := procLockWorkStation.Call()
	if r1 == 0 {
		return err
	}
	return nil
}
func getPlatformBatteryInfo() BatteryInfo {
	var sps SYSTEM_POWER_STATUS
	r1, _, _ := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&sps)))
	if r1 == 0 {
		return BatteryInfo{HasBattery: false, PowerSource: "AC / Unknown"}
	}

	hasBattery := sps.BatteryFlag != 128 && sps.BatteryFlag != 255
	pct := int(sps.BatteryLifePercent)
	if pct > 100 {
		pct = 100
	}

	powerSource := "Battery"
	if sps.ACLineStatus == 1 {
		powerSource = "AC Adapter (Plugged In)"
	}

	isCharging := (sps.BatteryFlag&8) != 0 || (sps.ACLineStatus == 1 && pct < 100)

	return BatteryInfo{
		HasBattery:  hasBattery,
		Percentage:  pct,
		IsCharging:  isCharging,
		PowerSource: powerSource,
	}
}

// SetMotorSpeed controls fan speeds from STOP (0) up to TURBO (5)
func SetMotorSpeed(keepMotor bool, motorLevel int) (uintptr, string) {
	var flags uintptr = ES_CONTINUOUS | ES_SYSTEM_REQUIRED
	var fanSpeed string

	if keepMotor {
		flags |= ES_DISPLAY_REQUIRED
		switch motorLevel {
		case 1:
			fanSpeed = "SPEED_1"
		case 2:
			fanSpeed = "SPEED_2"
		case 3:
			fanSpeed = "SPEED_3"
		case 4:
			fanSpeed = "SPEED_4"
		case 5:
			fanSpeed = "SPEED_5"
		case 0:
			fanSpeed = "STOP"
		default:
			fanSpeed = "SPEED_1"
		}
	}

	return flags, fanSpeed
}
