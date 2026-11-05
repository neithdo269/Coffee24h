// ==========================================
// PIECE 1: APPLICATION STATE & SAFEGUARD
// ==========================================
// Plain English: "The single source of truth: holds what the app is doing right now."
// Clean Code: 'const' guarantees the reference cannot be overwritten accidentally.
const state = {
    isAwake: false,
    mode: 'timed',
    keepDisplay: true,
    selectedMinutes: 15,
    secondsLeft: 15 * 60,
    targetSeconds: 15 * 60,
    battery: {
        hasBattery: false,
        percentage: 100,
        isCharging: false,
        powerSource: "AC Adapter (Plugged In)"
    },
    batteryGuard: true,
    safeguardPct: 5,
    autoLockOnEnd: false
};

// [OFFENSIVE DEFENSE: TAMPER-PROOF BATTERY SAFEGUARD]
// 1. Immutable Baseline: Frozen so runtime memory inspection cannot tamper with thresholds.
// 2. Failsafe Default: Provides safe telemetry if Win32 IPC drops.
const batterySafeGuard = Object.freeze({
    hasBattery: false,
    percentage: 100,
    isCharging: false,
    powerSource: "AC Adapter (Plugged In)",

    // Active defensive check:
    isCritical(currentPct) {
        return typeof currentPct === 'number' && currentPct <= 5;
    },
    sanitize(pct) {
        if (typeof pct !== 'number' || isNaN(pct)) return 100;
        return Math.max(0, Math.min(100, Math.floor(pct)));
    }
});

// ==========================================
// PIECE 2: TIME FORMATTER (Pure Function)
// ==========================================
// Plain English: "Convert raw seconds into a human clock (HH:MM:SS)"
// Clean Code Rule: A function should do ONE thing well with zero surprise side effects.
function formatTime(totalSeconds) {
    // Defense: If someone passes negative number, NaN, or string -> default safely
    if (typeof totalSeconds !== 'number' || isNaN(totalSeconds) || totalSeconds < 0) {
        return "00:00:00";
    }

    // Math breakdown:
    // 1 hour = 3600 seconds
    const hours = Math.floor(totalSeconds / 3600);
    // Remainder seconds turned into minutes (60 seconds = 1 minute)
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    // Leftover seconds
    const seconds = Math.floor(totalSeconds % 60);

    // "padStart(2, '0')" ensures 5 becomes "05", but 25 stays "25"
    const pad = (num) => String(num).padStart(2, '0');

    return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}

// ==========================================
// PIECE 3: DURATION PRESETS & STATE MUTATION
// ==========================================
// Plain English: "Quick 1-touch duration buttons (15m, 30m, 1h, 2h)."
// Engineer Jargon: "Duration Presets & State Mutation with Guard Clauses."
// Clean Code Rule: Guard against invalid inputs and isolate state changes cleanly.

const PRESET_MINUTES = Object.freeze([15, 30, 60, 120]);

function selectPreset(minutes) {
    // 1. Guard Clause: Don't change time while countdown is running
    if (state.isAwake) {
        console.warn("Session already running. Stop session before changing duration.");
        return;
    }

    // 2. Defensive Validation: Reject bad or negative values
    if (typeof minutes !== 'number' || isNaN(minutes) || minutes <= 0) {
        console.warn("Invalid preset minutes passed:", minutes);
        return;
    }

    // 3. Update single source of truth (state)
    state.selectedMinutes = Math.floor(minutes);
    state.targetSeconds = state.selectedMinutes * 60;
    state.secondsLeft = state.targetSeconds;

    // 4. Re-sync screen to show new target time
    updateUI();
}

// ==========================================
// PIECE 4: BATTERY TELEMETRY & BADGE UI
// ==========================================
// Plain English: "Show live battery percentage, charger status, and warn if battery is low."
// Engineer Jargon: "Hardware Telemetry Mapping & State-Driven Badge Renderer."
// Clean Code: Pure helper maps battery data -> visual presentation.

function getBatteryVisuals(battery) {
    // 1. Máy bàn (PC) không có pin, cắm nguồn trực tiếp
    if (!battery || !battery.hasBattery) {
        return {
            icon: "⚡",
            text: "AC Power • Plugged In",
            badgeClass: "charging"
        };
    }

    // 2. Làm sạch số % pin qua chốt bảo vệ Piece 1
    const pct = batterySafeGuard.sanitize(battery.percentage);

    // 3. Đang cắm sạc
    if (battery.isCharging) {
        return {
            icon: "⚡",
            text: `${pct}% • Charging`,
            badgeClass: "charging"
        };
    }

    // 4. Đang dùng pin - cảnh báo nếu pin chạm mức nguy hiểm (<= 5%)
    if (batterySafeGuard.isCritical(pct)) {
        return {
            icon: "🪫",
            text: `${pct}% • Critical Battery!`,
            badgeClass: "critical"
        };
    }

    // 5. Đang dùng pin bình thường
    return {
        icon: "🔋",
        text: `${pct}% • On Battery`,
        badgeClass: "normal"
    };
}

// Cập nhật huy hiệu pin trên màn hình
function updateBatteryUI() {
    const badgeEl = document.getElementById("batteryBadge");
    const iconEl = document.getElementById("batteryIcon");
    const textEl = document.getElementById("batteryText");

    if (!badgeEl || !iconEl || !textEl) return;

    const visuals = getBatteryVisuals(state.battery);

    iconEl.textContent = visuals.icon;
    textEl.textContent = visuals.text;
    badgeEl.className = `battery-badge ${visuals.badgeClass}`;
}

// ==========================================
// PIECE 5: THE TIMER CONTROLLER & UI SYNC
// ==========================================
// Plain English: "Start ticking down the clock when ON, stop when OFF, and update what the user sees."

let timerInterval = null; // Holds the active ticker

// Update the screen so it always matches our 'state' object
function updateUI() {
    const secsEl = document.getElementById("secs");
    const toggleBtn = document.getElementById("toggleBtn");
    const presetButtons = document.querySelectorAll(".preset-chip");
    const presetsContainer = document.getElementById("presetsContainer");
    const modeButtons = document.querySelectorAll(".mode-chip");

    // 1. Update Battery Status Badge
    updateBatteryUI();

    // 2. Update Mode Chips
    modeButtons.forEach((btn) => {
        btn.classList.toggle("active", btn.dataset.mode === state.mode);
        btn.disabled = state.isAwake;
    });

    // 3. Show presets only when in 'timed' mode
    if (presetsContainer) {
        presetsContainer.style.display = state.mode === "timed" ? "flex" : "none";
    }

    // 4. Update Timer Display
    if (secsEl) {
        if (state.mode === "until_off") {
            secsEl.textContent = state.isAwake ? "♾️ AWAKE" : "♾️ UNTIL OFF";
        } else if (state.mode === "system_only") {
            secsEl.textContent = state.isAwake ? "💻 BG RUN" : "💻 SYSTEM ONLY";
        } else {
            const displaySecs = state.isAwake ? state.secondsLeft : state.targetSeconds;
            secsEl.textContent = formatTime(displaySecs);
        }
    }

    // 5. Update Preset Chips: Highlight active chip, lock chips while running
    presetButtons.forEach((btn) => {
        const mins = parseInt(btn.dataset.minutes, 10);
        const isCurrent = mins === state.selectedMinutes;

        btn.classList.toggle("active", isCurrent);
        btn.disabled = state.isAwake; // Freeze selection during active countdown
    });

    // 6. Update Big Toggle Button
    if (toggleBtn) {
        if (state.isAwake) {
            toggleBtn.textContent = "Stop Awake";
            toggleBtn.classList.add("active");
        } else {
            toggleBtn.textContent = "Start Awake";
            toggleBtn.classList.remove("active");
        }
    }
}

// Start counting down every 1 second
function startCountdown(durationSeconds) {
    stopCountdown(); // Always clear previous timer before starting a new one

    state.isAwake = true;
    state.secondsLeft = durationSeconds;
    updateUI();

    // 1-second heartbeat
    timerInterval = setInterval(() => {
        if (state.secondsLeft > 0) {
            state.secondsLeft--;
            updateUI();
        } else {
            // Reached 00:00:00 -> Session finished
            stopCountdown();
        }
    }, 1000);
}

// Stop counting down and reset
function stopCountdown() {
    if (timerInterval) {
        clearInterval(timerInterval);
        timerInterval = null;
    }
    state.isAwake = false;
    state.secondsLeft = state.targetSeconds; // Reset clock back to selected preset
    updateUI();
}

// ==========================================
// PIECE 7: AUDIO FEEDBACK (Pure Web Audio API)
// ==========================================
// Plain English: "Pleasant chime sound when coffee starts, calming sound when stopping."
// Engineer Jargon: "Offline Synthesized Acoustic Feedback (Zero External Assets)."
function playChime(isStart) {
    try {
        const AudioContext = window.AudioContext || window.webkitAudioContext;
        if (!AudioContext) return;
        const ctx = new AudioContext();

        const freqs = isStart ? [523.25, 659.25, 783.99] : [659.25, 523.25, 392.00];
        freqs.forEach((freq, idx) => {
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();

            osc.type = isStart ? 'triangle' : 'sine';
            const startTime = ctx.currentTime + (idx * 0.06);
            osc.frequency.setValueAtTime(freq, startTime);

            gain.gain.setValueAtTime(0.12, startTime);
            gain.gain.exponentialRampToValueAtTime(0.001, startTime + 0.35);

            osc.connect(gain);
            gain.connect(ctx.destination);
            osc.start(startTime);
            osc.stop(startTime + 0.4);
        });
    } catch (err) {
        // Audio synthesis safely ignored if user hasn't interacted with page yet
    }
}

// ==========================================
// PIECE 8: MODE SELECTION & CONTROLS
// ==========================================
// Plain English: "Change between Until Off, Timed, or System Only modes."
function selectMode(newMode) {
    if (state.isAwake) {
        console.warn("Stop session before changing mode.");
        return;
    }
    state.mode = newMode;
    updateUI();
}

// When user clicks the big button
async function handleToggleClick() {
    if (state.isAwake) {
        // Stop awake session
        stopCountdown();
        playChime(false);
        if (window?.go?.main?.App?.StopAwake) {
            await window.go.main.App.StopAwake();
        }
    } else {
        // Start awake session
        playChime(true);
        if (window?.go?.main?.App?.StartAwake) {
            state.isAwake = true;
            updateUI();
            await window.go.main.App.StartAwake(
                state.mode,
                state.keepDisplay,
                state.selectedMinutes,
                state.batteryGuard,
                state.safeguardPct,
                state.autoLockOnEnd
            );
        } else {
            // Standalone browser fallback without Wails
            startCountdown(state.targetSeconds);
        }
    }
}

// ==========================================
// PIECE 6 & 9: SAFE WAILS BRIDGE & INITIALIZER
// ==========================================
// Plain English: "Wire up buttons, checkboxes, and talk to Go backend if inside desktop app."
async function init() {
    // 1. Hook up the big button click event
    const toggleBtn = document.getElementById("toggleBtn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", handleToggleClick);
    }

    // 2. Hook up Preset Chips
    const presetButtons = document.querySelectorAll(".preset-chip");
    presetButtons.forEach((btn) => {
        btn.addEventListener("click", () => {
            const mins = parseInt(btn.dataset.minutes, 10);
            selectPreset(mins);
        });
    });

    // 3. Hook up Mode Chips (Until Off, Timed, System Only)
    const modeButtons = document.querySelectorAll(".mode-chip");
    modeButtons.forEach((btn) => {
        btn.addEventListener("click", () => {
            selectMode(btn.dataset.mode);
        });
    });

    // 4. Hook up Safeguard Checkboxes
    const chkKeepDisplay = document.getElementById("chkKeepDisplay");
    if (chkKeepDisplay) {
        chkKeepDisplay.addEventListener("change", (e) => {
            state.keepDisplay = e.target.checked;
        });
    }

    const chkSafeguard = document.getElementById("chkSafeguard");
    if (chkSafeguard) {
        chkSafeguard.addEventListener("change", (e) => {
            state.batteryGuard = e.target.checked;
        });
    }

    const chkAutoLock = document.getElementById("chkAutoLock");
    if (chkAutoLock) {
        chkAutoLock.addEventListener("change", (e) => {
            state.autoLockOnEnd = e.target.checked;
        });
    }

    // 5. Hook up Instant Lock Workstation Button
    const btnLockNow = document.getElementById("btnLockNow");
    if (btnLockNow) {
        btnLockNow.addEventListener("click", async () => {
            if (window?.go?.main?.App?.LockNow) {
                await window.go.main.App.LockNow();
            } else if (window?.go?.main?.App?.LockWorkstation) {
                await window.go.main.App.LockWorkstation();
            } else {
                console.log("Mock Lock Workstation triggered (Win + L).");
            }
        });
    }

    // 6. Fetch initial calculation and state from backend if available
    let initialSecs = state.targetSeconds;
    if (window?.go?.main?.App?.CalculateSeconds) {
        initialSecs = await window.go.main.App.CalculateSeconds(state.selectedMinutes);
        state.targetSeconds = initialSecs;
        state.secondsLeft = initialSecs;
    }

    if (window?.go?.main?.App?.GetState) {
        const backendState = await window.go.main.App.GetState();
        if (backendState) {
            if (backendState.battery) {
                state.battery = backendState.battery;
            }
            if (typeof backendState.isAwake === 'boolean') {
                state.isAwake = backendState.isAwake;
            }
        }
    }

    // 7. Listen for live hardware battery updates & timer ticks from Go backend
    if (window?.runtime?.EventsOn) {
        window.runtime.EventsOn("state_updated", (backendState) => {
            if (backendState) {
                if (backendState.battery) {
                    state.battery = backendState.battery;
                    updateBatteryUI();
                }
                if (typeof backendState.isAwake === 'boolean') {
                    state.isAwake = backendState.isAwake;
                }
                if (typeof backendState.secondsLeft === 'number' && state.mode === 'timed') {
                    state.secondsLeft = backendState.secondsLeft;
                }
                updateUI();
            }
        });

        // Listen for safeguard trigger (pin yếu tự tắt)
        window.runtime.EventsOn("safeguard_triggered", (msg) => {
            console.warn("SAFEGUARD:", msg);
            stopCountdown();
        });

        // Listen for timer completion
        window.runtime.EventsOn("timer_finished", (msg) => {
            console.log("TIMER FINISHED:", msg);
            stopCountdown();
        });
    }

    console.log(`Coffee24h initialized! Mode: ${state.mode}, Duration: ${state.selectedMinutes}m`);
    updateUI();
}

// Start when document is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}


