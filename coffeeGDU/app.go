package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetGduInfo() string {
	return "GDU Window Manager - Coffee24h v1.0 by Neith Do"
}

func (a *App) CalculateSeconds(minutes int) int {
	return minutes * 60
}

// SetCupEdition switches between "gdu" (Classic Blue Package Cup) and "eco" (Student Recycle Leaf Cup)
func (a *App) SetCupEdition(edition string) string {
	switchCupIcon(edition)
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "cup_edition_changed", edition)
	}
	return edition
}

// GetCupEdition returns active cup edition
func (a *App) GetCupEdition() string {
	if a.cliOpts.EcoCup {
		return "eco"
	}
	return currentCupEdition
}

// AwakeState is passed to the frontend
type AwakeState struct {
	IsAwake       bool        `json:"isAwake"`
	Mode          string      `json:"mode"` // "until_off", "timed", "system_only"
	KeepDisplay   bool        `json:"keepDisplay"`
	SecondsLeft   int         `json:"secondsLeft"`   // 0 if until_off
	TargetSeconds int         `json:"targetSeconds"` // total duration if timed
	Battery       BatteryInfo `json:"battery"`
	BatteryGuard  bool        `json:"batteryGuard"` // stop if battery < 5%
	SafeguardPct  int         `json:"safeguardPct"`
	AutoLockOnEnd bool        `json:"autoLockOnEnd"` // auto-lock Windows when session ends
}

// App struct
type App struct {
	ctx        context.Context
	mu         sync.Mutex // The recommended Mutex in Go is sync.Mutex.
	// Called "THE SHIELD". the Goroutine A & B "Stop Awake"&  00:00:00. At the exact
	// same time "Millisecond" the roundabout Race condition / Memory Corruption in Go will lead to undefined behavior
	// App crashed at this point of without this "my sync.Mutex"
	inhibitor  PlatformInhibitor
	state      AwakeState
	timerDone  chan struct{}
	tickerStop chan struct{}
	cliOpts    CLIOptions
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		inhibitor: newPlatformInhibitor(),
		state: AwakeState{
			IsAwake:       false,
			Mode:          "until_off",
			KeepDisplay:   true,
			BatteryGuard:  true,
			SafeguardPct:  5,
			AutoLockOnEnd: false,
		},
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.state.Battery = getPlatformBatteryInfo()

	// Initialize system tray icon in "Show hidden icons"
	initTray(a)

	// Start background poller for battery & status tick
	go a.statusLoop()

	// Process CLI startup flags
	if a.cliOpts.AwakeNow || a.cliOpts.TimedMin > 0 || a.cliOpts.SystemOnly {
		mode := "until_off"
		keepDisplay := true
		if a.cliOpts.SystemOnly {
			mode = "system_only"
			keepDisplay = false
		} else if a.cliOpts.TimedMin > 0 {
			mode = "timed"
		}
		safeguard := !a.cliOpts.NoSafeguard
		go func() {
			time.Sleep(300 * time.Millisecond)
			_, _ = a.StartAwake(mode, keepDisplay, a.cliOpts.TimedMin, safeguard, 5, a.cliOpts.AutoLock)
		}()
	}
}

// domReady anchors the window at the bottom-right corner of the primary screen above the taskbar.
func (a *App) domReady(ctx context.Context) {
	if a.cliOpts.EcoCup {
		switchCupIcon("eco")
		runtime.EventsEmit(ctx, "cup_edition_changed", "eco")
	} else {
		setWindowIcons("Coffee24h ☕ GDU")
	}
	screens, err := runtime.ScreenGetAll(ctx)
	if err == nil && len(screens) > 0 {
		primary := screens[0]
		for _, s := range screens {
			if s.IsPrimary {
				primary = s
				break
			}
		}

		const (
			winW   = 350
			winH   = 350
			margin = 16
			tbH    = 48 // Standard Windows taskbar height offset
		)

		x := primary.Width - winW - margin
		y := primary.Height - winH - margin - tbH
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}

		runtime.WindowSetSize(ctx, winW, winH)
		runtime.WindowSetPosition(ctx, x, y)
	}
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_ = a.inhibitor.Release()
	a.stopTimerLocked()
	removeTray()
}

// Quit completely shuts down the application from the tray context menu
func (a *App) Quit() {
	a.mu.Lock()
	_ = a.inhibitor.Release()
	a.stopTimerLocked()
	a.mu.Unlock()
	removeTray()
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

func (a *App) statusLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.mu.Lock()
			a.state.Battery = getPlatformBatteryInfo()
			// [SECURITY: STAGE-FAIL-SAFE & THEFT DEFENSE]
			//* 1. Alert State: Audible feedback alerts the user even if the monitor is dark.
			//* 2. Physical Security: Auto-locking prevents "evil maid" / walk-by theft if
			//     the laptop is abandoned in a public cafe and battery drains to critical.
			//* 3. Portability: Requires user32.dll (LockWorkStation & MessageBeep).
			//* 4. Safeguard Trigger: If awake on battery and percentage drops below threshold.
			if a.state.IsAwake && a.state.BatteryGuard && a.state.Battery.HasBattery && !a.state.Battery.IsCharging {
				if a.state.Battery.Percentage <= a.state.SafeguardPct {
					_ = a.inhibitor.Release()
					a.state.IsAwake = false
					a.stopTimerLocked()
					updateTrayStatus(false)
					PlaySystemSound(0x10) // Critical Stop / Alert sound
					// SECURITY PATTERN: Zero-Trust Auto-Lock & Telemetry Alert
					// Prevents unauthorized physical access when unattended battery depletion occurs.
					if a.state.AutoLockOnEnd {
						_ = LockScreen()
					}
					runtime.EventsEmit(a.ctx, "safeguard_triggered", fmt.Sprintf("Battery reached %d%%! Awake deactivated to protect work.", a.state.Battery.Percentage))
				}
			}

			// If window was minimized via OS shortcut, hide it to system tray
			if a.ctx != nil && runtime.WindowIsMinimised(a.ctx) {
				runtime.WindowHide(a.ctx)
			}

			runtime.EventsEmit(a.ctx, "state_updated", a.state)
			a.mu.Unlock()
		}
	}
}

// GetState returns the current awake state
func (a *App) GetState() AwakeState {
	a.mu.Lock() //Lock "Tn" the mutex to prevent race condition
	defer a.mu.Unlock()//Unlock "Fn" the mutex to prevent race condition
	a.state.Battery = getPlatformBatteryInfo()
	return a.state
}

// [SECURITY: ZERO-TRUST INPUT DEFENSE & SANITIZATION]
// * 1. Boundary Defense: Never trust caller arguments from Webview / IPC.
// * 2. DoS Prevention: Clamps duration (0 to 1440 min / 24 hours max) to prevent overflow.
// * 3. Battery Safeguard Bound: Clamps safeguard percentage (3% to 50%) to prevent critical bricking.
// * 4. Thread-Safe Transition: Acquires mutex lock before any state modifications.
func (a *App) StartAwake(mode string, keepDisplay bool, durationMinutes int, batteryGuard bool, safeguardPct int, autoLock bool) (AwakeState, error) {
	// [SECURITY: SYNCHRONIZATION]
	// [MOCK QUESTION]
	// Question: "When goroutine A try to turn on awake and at the same time goroutine B try to turn off awake, what will happen ?"
	// Answer: "Race condition / Memory Corruption in Go will lead to undefined behavior"
	// a.mu.Lock is the "THE THIRD SHIELD" that prevent this from happening
	a.mu.Lock() // Lock "Tn" the mutex to prevent race condition
	defer a.mu.Unlock()//Unlock "Fn" the mutex to prevent race condition

	// Zero-Trust input defense: sanitize and clamp values
	if durationMinutes < 0 {
		durationMinutes = 0
	}
	if durationMinutes > 1440 { // Max 24 hours
		durationMinutes = 1440
	}
	if safeguardPct < 3 {
		safeguardPct = 3
	} else if safeguardPct > 50 {
		safeguardPct = 50
	}

	a.stopTimerLocked()

	err := a.inhibitor.Inhibit(keepDisplay)
	if err != nil {
		return a.state, err
	}

	a.state.IsAwake = true
	a.state.Mode = mode
	a.state.KeepDisplay = keepDisplay
	a.state.BatteryGuard = batteryGuard
	a.state.SafeguardPct = safeguardPct
	a.state.AutoLockOnEnd = autoLock
	if mode == "timed" && durationMinutes > 0 {
		a.state.TargetSeconds = durationMinutes * 60
		a.state.SecondsLeft = a.state.TargetSeconds
		a.startTimerLocked()
	} else {
		a.state.TargetSeconds = 0
		a.state.SecondsLeft = 0
	}

	a.state.Battery = getPlatformBatteryInfo()
	updateTrayStatus(true)
	runtime.EventsEmit(a.ctx, "state_updated", a.state)
	return a.state, nil
}

// StopAwake stops the keep-awake session
func (a *App) StopAwake() AwakeState {
	a.mu.Lock()         // 🔒 LOCK "Tn" the mutex to prevent race condition
	defer a.mu.Unlock() // 🔓 UNLOCK "Fn" the mutex when function returns

	_ = a.inhibitor.Release()
	a.stopTimerLocked()
	a.state.IsAwake = false
	a.state.SecondsLeft = 0
	a.state.TargetSeconds = 0

	a.state.Battery = getPlatformBatteryInfo()
	updateTrayStatus(false)
	runtime.EventsEmit(a.ctx, "state_updated", a.state)
	return a.state
}

// ToggleAwake toggles awake state with default until-off settings
func (a *App) ToggleAwake() AwakeState {
	a.mu.Lock()
	isCurrentlyAwake := a.state.IsAwake
	a.mu.Unlock()

	if isCurrentlyAwake {
		return a.StopAwake()
	}
	st, _ := a.StartAwake("until_off", true, 0, true, 5, false)
	return st
}

func (a *App) startTimerLocked() {
	timerStop := make(chan struct{})
	a.timerDone = timerStop

	go func() {
		tick := time.NewTicker(1 * time.Second)
		defer tick.Stop()

		for {
			select {
			case <-timerStop:
				return
			case <-tick.C:
				a.mu.Lock()
				// [SECURITY: SESSION EXPIRATION & WORKSTATION SHIELD]
				//* 1. Session Teardown: Releases Win32 thread execution state.
				//* 2. Audio Telemetry: 0x40 (MB_ICONASTERISK) produces a discreet completion chime.
				//* 3. Zero-Trust Lockdown: Forces Windows into lock screen (Win + L) if AutoLock is enabled.
				//* 4. Frontend Dispatch: Emits real-time event to update the UI countdown banner.
				if !a.state.IsAwake || a.state.SecondsLeft <= 0 {
					_ = a.inhibitor.Release()
					a.state.IsAwake = false
					a.state.SecondsLeft = 0
					updateTrayStatus(false)
					PlaySystemSound(0x40) // Information / Chime sound
					if a.state.AutoLockOnEnd {
						_ = LockScreen()
					}
					runtime.EventsEmit(a.ctx, "state_updated", a.state)
					runtime.EventsEmit(a.ctx, "timer_finished", "Stay awake session timed out.")
					a.mu.Unlock()
					return
				}
				a.state.SecondsLeft--
				runtime.EventsEmit(a.ctx, "state_updated", a.state)
				a.mu.Unlock()
			}
		}
	}()
}

func (a *App) stopTimerLocked() {
	if a.timerDone != nil {
		close(a.timerDone)
		a.timerDone = nil
	}
}

// LockNow immediately locks the workstation on demand
func (a *App) LockNow() error {
	return LockScreen()
}

// LockWorkstation is an alias for LockNow to maintain full compatibility across frontend variants
func (a *App) LockWorkstation() error {
	return LockScreen()
}

// MinimizeWindow hides the application window to the system tray ("Show hidden icons")
func (a *App) MinimizeWindow() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

// HideWindow hides the application window to the system tray
func (a *App) HideWindow() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

