package main

// PlatformInhibitor interface for cross-platform OS power management
type PlatformInhibitor interface {
	Inhibit(keepDisplay bool) error
	Release() error
	IsActive() bool
}

// BatteryInfo holds battery statistics across OSes
type BatteryInfo struct {
	HasBattery bool    `json:"hasBattery"`
	Percentage int     `json:"percentage"`
	IsCharging bool    `json:"isCharging"`
	PowerSource string `json:"powerSource"`
}
