import './style.css';
import coffeeLogo from './assets/images/coffee-logo.png';

// Bilingual Dictionary (Tiếng Việt & English)
const i18n = {
  vi: {
    appTitle: "Coffee24h",
    appSubtitle: "Đại học Gia Định (GDU) - Bật Ý Tưởng",
    statusAwakeUntilOff: "Bật Ý Tưởng: Chạy Không Giới Hạn",
    statusAwakeTimed: "Bật Ý Tưởng: Đang Hẹn Giờ",
    statusAwakeSystemOnly: "Chạy Nền: Tắt Màn Hình Tiết Kiệm Pin",
    statusSleepAllowed: "Cho Phép Nghỉ Ngơi (Sleep)",
    descUntilOff: "Máy tính luôn thức sáng suốt đêm cho đến khi bạn tắt hoặc hết pin",
    descTimed: "Đang đếm ngược thời gian tập trung...",
    descSystemOnly: "Màn hình có thể tắt; CPU và tác vụ tải/render vẫn chạy liên tục",
    descSleep: "Máy tính sẽ ngủ theo chính sách nguồn của Windows",
    modeUntilOff: "Vô Hạn",
    modeTimed: "Hẹn Giờ",
    modeSystemOnly: "Chạy Nền",
    sectionMode: "Chế Độ Bật Sáng",
    keepDisplay: "Giữ Màn Hình Sáng",
    keepDisplayDesc: "Chống tắt màn hình và ngăn màn hình chờ (screensaver)",
    safeguard: "Bảo Vệ Pin Thông Minh (5%)",
    safeguardDesc: "Tự động tắt keep-awake khi pin < 5% để không sập nguồn đột ngột",
    autoLock: "🛡️ Tự Động Khóa Màn Hình (Win + L)",
    autoLockDesc: "Tự khóa máy bảo mật khi hết giờ hoặc khi pin yếu rời bàn làm việc",
    lockNowBtn: "🔒 Khóa Máy Ngay Lập Tức",
    footerNotice: "Tương thích hệ điều hành Windows Win32 API • Made with ❤️ for GDU",
    acPlugged: "Đang cắm sạc (AC)",
    batteryPower: "Đang dùng Pin"
  },
  en: {
    appTitle: "Coffee24h",
    appSubtitle: "Cross-Platform Keep-Awake",
    statusAwakeUntilOff: "Stay Awake: Until Power Off",
    statusAwakeTimed: "Stay Awake: Timed Session",
    statusAwakeSystemOnly: "Stay Awake: System Only",
    statusSleepAllowed: "Sleep Allowed",
    descUntilOff: "Laptop will stay powered on until battery exhausts or manual stop",
    descTimed: "Active focus countdown in progress...",
    descSystemOnly: "Display can turn off; CPU and background tasks stay awake",
    descSleep: "Laptop will sleep according to OS power policy",
    modeUntilOff: "Until Off",
    modeTimed: "Timed",
    modeSystemOnly: "System Only",
    sectionMode: "Awake Mode",
    keepDisplay: "Keep Display On",
    keepDisplayDesc: "Prevent screen dimming and screensaver",
    safeguard: "Battery Safeguard (5%)",
    safeguardDesc: "Auto-release awake mode before laptop completely shuts off",
    autoLock: "🛡️ Auto-Lock on Session End",
    autoLockDesc: "Lock Windows screen (Win + L) when timer ends or battery drops",
    lockNowBtn: "🔒 Lock Workstation Now",
    footerNotice: "Compatible with Windows (SetThreadExecutionState), macOS & Linux",
    acPlugged: "AC Adapter (Plugged In)",
    batteryPower: "On Battery"
  }
};

let currentLang = 'vi'; // Default to Vietnamese for the authentic Café experience!

// Local UI state
let state = {
  isAwake: false,
  mode: 'until_off', // 'until_off', 'timed', 'system_only'
  keepDisplay: true,
  secondsLeft: 0,
  targetSeconds: 0,
  batteryGuard: true,
  safeguardPct: 5,
  autoLockOnEnd: false,
  battery: {
    hasBattery: false,
    percentage: 100,
    isCharging: false,
    powerSource: 'Detecting...'
  }
};

let selectedMinutes = 30;

// [ACOUSTICS: BẬT Ý TƯỞNG AUDIO SYNTHESIZER]
// 100% offline Web Audio API harmonic sound - like a ceramic coffee cup chime igniting late-night focus
function playBatyTuongSound() {
  try {
    const AudioCtx = window.AudioContext || window.webkitAudioContext;
    if (!AudioCtx) return;
    const ctx = new AudioCtx();

    // Harmonic notes: C5 (523Hz), E5 (659Hz), G5 (784Hz), C6 (1046Hz) shimmer
    const chord = [523.25, 659.25, 783.99, 1046.50];
    chord.forEach((freq, idx) => {
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'triangle'; // Warm, bell-like timbre
      const startTime = ctx.currentTime + idx * 0.055;
      osc.frequency.setValueAtTime(freq, startTime);

      gain.gain.setValueAtTime(0.16, startTime);
      gain.gain.exponentialRampToValueAtTime(0.001, startTime + 0.45);

      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start(startTime);
      osc.stop(startTime + 0.5);
    });
  } catch (e) {
    console.warn('Audio synthesis not permitted or blocked:', e);
  }
}

// Gentle Power-Down Sound when entering Sleep Mode
function playSleepRestSound() {
  try {
    const AudioCtx = window.AudioContext || window.webkitAudioContext;
    if (!AudioCtx) return;
    const ctx = new AudioCtx();

    const notes = [659.25, 523.25, 392.00]; // E5 -> C5 -> G4 calming descent
    notes.forEach((freq, idx) => {
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'sine';
      const startTime = ctx.currentTime + idx * 0.07;
      osc.frequency.setValueAtTime(freq, startTime);

      gain.gain.setValueAtTime(0.12, startTime);
      gain.gain.exponentialRampToValueAtTime(0.001, startTime + 0.35);

      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start(startTime);
      osc.stop(startTime + 0.4);
    });
  } catch (e) {
    // Silent fallback
  }
}

// Setup App HTML
const app = document.querySelector('#app');
app.innerHTML = `
  <div class="app-container">
    <!-- Header: Clean GDU Transparent Glass Title Bar (No text overlap, No side bar) -->
    <header class="app-header">
      <div class="brand">
        <img src="${coffeeLogo}" alt="Coffee24h" style="width: 24px; height: 24px; border-radius: 4px; object-fit: contain; margin-right: 6px;" />
        <h1 class="brand-title"><span id="txtTitle">Coffee24h</span><span class="gdu-tag">GDU</span></h1>
        <span id="txtSubtitle" style="display: none;"></span>
      </div>
      <div class="window-controls">
        <button id="btnWinMin" class="win-btn win-min" title="Thu nhỏ vào khay hệ thống / Minimize to tray">−</button>
        <button id="btnWinClose" class="win-btn win-close" title="Đóng vào khay / Close to tray">✕</button>
      </div>
    </header>

    <!-- Hero Card (Khung Bật Ý Tưởng) -->
    <section id="heroCard" class="hero-toggle-card">
      <div class="hero-header-row">
        <button id="btnLang" class="lang-switch-btn" title="Chuyển đổi ngôn ngữ / Switch Language">
          vn-vi
        </button>
        <div id="batteryBadge" class="battery-badge">
          <span id="batIcon">⚡</span>
          <span id="batPct">--%</span>
        </div>
      </div>

      <div class="steam-container">
        <div class="steam-particle"></div>
        <div class="steam-particle"></div>
        <div class="steam-particle"></div>
      </div>
      <div class="steam-ring"></div>
      <button id="mainToggleBtn" class="main-toggle-btn" title="Click để Bật Ý Tưởng / Toggle Awake">
        <span class="cup-icon">☕</span>
      </button>

      <div id="statusHeadline" class="status-headline">Sleep Allowed</div>
      <div id="statusSubtext" class="status-subtext">Laptop will sleep according to OS power policy</div>
      <div id="countdownBanner" class="countdown-banner">00:00:00</div>
    </section>

    <!-- Modes Selection (Khung Chế Độ) -->
    <div class="mode-card">
      <div class="section-title" id="txtSectionMode">Chế Độ Bật Sáng</div>
      <div class="mode-grid">
        <button id="modeUntilOff" class="mode-btn active" data-mode="until_off">
          <span class="icon">♾️</span>
          <span class="label" id="lblModeUntilOff">Vô Hạn</span>
        </button>
        <button id="modeTimed" class="mode-btn" data-mode="timed">
          <span class="icon">⏱️</span>
          <span class="label" id="lblModeTimed">Hẹn Giờ</span>
        </button>
        <button id="modeSystemOnly" class="mode-btn" data-mode="system_only">
          <span class="icon">💻</span>
          <span class="label" id="lblModeSystemOnly">Chạy Nền</span>
        </button>
      </div>

      <!-- Presets for Timed Mode -->
      <div id="presetsBar" class="presets-bar" style="display: none;">
        <button class="preset-chip" data-min="15">15m</button>
        <button class="preset-chip active" data-min="30">30m</button>
        <button class="preset-chip" data-min="60">1h</button>
        <button class="preset-chip" data-min="120">2h</button>
        <button class="preset-chip" data-min="240">4h</button>
      </div>
    </div>

    <!-- Toggles & Safeguards (Khung Cài Đặt Chính) -->
    <div class="controls-card">
      <div class="toggle-row">
        <div class="toggle-info">
          <span class="toggle-label" id="lblKeepDisplay">Giữ Màn Hình Sáng</span>
          <span class="toggle-desc" id="descKeepDisplay">Chống tắt màn hình và ngăn màn hình chờ</span>
        </div>
        <label class="switch">
          <input type="checkbox" id="chkKeepDisplay" checked />
          <span class="slider"></span>
        </label>
      </div>

      <div class="toggle-row">
        <div class="toggle-info">
          <span class="toggle-label" id="lblSafeguard">Bảo Vệ Pin Thông Minh (5%)</span>
          <span class="toggle-desc" id="descSafeguard">Tự động tắt awake khi pin < 5% để không sập nguồn</span>
        </div>
        <label class="switch">
          <input type="checkbox" id="chkSafeguard" checked />
          <span class="slider"></span>
        </label>
      </div>

      <div class="toggle-row">
        <div class="toggle-info">
          <span class="toggle-label" id="lblAutoLock">🛡️ Tự Động Khóa Màn Hình (Win + L)</span>
          <span class="toggle-desc" id="descAutoLock">Khóa máy an toàn khi hết giờ hoặc khi pin yếu</span>
        </div>
        <label class="switch">
          <input type="checkbox" id="chkAutoLock" />
          <span class="slider"></span>
        </label>
      </div>

      <button id="btnLockNow" class="lock-now-btn" title="Khóa ngay máy tính (Win + L)">
        🔒 Khóa Máy Ngay Lập Tức
      </button>

      <!-- Dòng chữ bản quyền & hệ điều hành nằm gọn gàng bên trong khung chính -->
      <div id="toastNotice" class="toast-notice">
        Tương thích hệ điều hành Windows Win32 API • Made with ❤️ for GDU
      </div>
    </div>
  </div>
  </div>
`;

// DOM Elements
const heroCard = document.getElementById('heroCard');
const mainToggleBtn = document.getElementById('mainToggleBtn');
const statusHeadline = document.getElementById('statusHeadline');
const statusSubtext = document.getElementById('statusSubtext');
const countdownBanner = document.getElementById('countdownBanner');
const batteryBadge = document.getElementById('batteryBadge');
const batIcon = document.getElementById('batIcon');
const batPct = document.getElementById('batPct');
const presetsBar = document.getElementById('presetsBar');
const chkKeepDisplay = document.getElementById('chkKeepDisplay');
const chkSafeguard = document.getElementById('chkSafeguard');
const chkAutoLock = document.getElementById('chkAutoLock');
const btnLockNow = document.getElementById('btnLockNow');
const btnLang = document.getElementById('btnLang');
const btnWinMin = document.getElementById('btnWinMin');
const btnWinClose = document.getElementById('btnWinClose');
const toastNotice = document.getElementById('toastNotice');

const txtTitle = document.getElementById('txtTitle');
const txtSubtitle = document.getElementById('txtSubtitle');
const txtSectionMode = document.getElementById('txtSectionMode');
const lblModeUntilOff = document.getElementById('lblModeUntilOff');
const lblModeTimed = document.getElementById('lblModeTimed');
const lblModeSystemOnly = document.getElementById('lblModeSystemOnly');
const lblKeepDisplay = document.getElementById('lblKeepDisplay');
const descKeepDisplay = document.getElementById('descKeepDisplay');
const lblSafeguard = document.getElementById('lblSafeguard');
const descSafeguard = document.getElementById('descSafeguard');
const lblAutoLock = document.getElementById('lblAutoLock');
const descAutoLock = document.getElementById('descAutoLock');

const modeButtons = {
  until_off: document.getElementById('modeUntilOff'),
  timed: document.getElementById('modeTimed'),
  system_only: document.getElementById('modeSystemOnly'),
};

function formatTime(seconds) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
}

function updateLanguageTexts() {
  const t = i18n[currentLang];
  btnLang.textContent = currentLang === 'vi' ? 'vn 🇻🇳-vi' : 'us en-us';
  txtTitle.textContent = t.appTitle;
  txtSubtitle.textContent = t.appSubtitle;
  txtSectionMode.textContent = t.sectionMode;
  lblModeUntilOff.textContent = t.modeUntilOff;
  lblModeTimed.textContent = t.modeTimed;
  lblModeSystemOnly.textContent = t.modeSystemOnly;
  lblKeepDisplay.textContent = t.keepDisplay;
  descKeepDisplay.textContent = t.keepDisplayDesc;
  lblSafeguard.textContent = t.safeguard;
  descSafeguard.textContent = t.safeguardDesc;
  lblAutoLock.textContent = t.autoLock;
  descAutoLock.textContent = t.autoLockDesc;
  btnLockNow.innerHTML = t.lockNowBtn;
  toastNotice.textContent = t.footerNotice;
}

function renderUI() {
  const t = i18n[currentLang];

  if (state.isAwake) {
    heroCard.classList.add('awake');
    statusHeadline.textContent = state.mode === 'until_off' ? t.statusAwakeUntilOff :
      state.mode === 'timed' ? t.statusAwakeTimed : t.statusAwakeSystemOnly;

    if (state.mode === 'until_off') {
      statusSubtext.textContent = t.descUntilOff;
      countdownBanner.classList.remove('visible');
    } else if (state.mode === 'timed') {
      statusSubtext.textContent = t.descTimed;
      countdownBanner.textContent = formatTime(state.secondsLeft);
      countdownBanner.classList.add('visible');
    } else {
      statusSubtext.textContent = t.descSystemOnly;
      countdownBanner.classList.remove('visible');
    }
  } else {
    heroCard.classList.remove('awake');
    statusHeadline.textContent = t.statusSleepAllowed;
    statusSubtext.textContent = t.descSleep;
    countdownBanner.classList.remove('visible');
  }

  // Active mode buttons
  Object.keys(modeButtons).forEach((m) => {
    if (m === state.mode) {
      modeButtons[m].classList.add('active');
    } else {
      modeButtons[m].classList.remove('active');
    }
  });

  presetsBar.style.display = state.mode === 'timed' ? 'flex' : 'none';

  // Battery status
  if (state.battery) {
    if (state.battery.hasBattery) {
      batPct.textContent = state.battery.percentage + '%';
      batIcon.textContent = state.battery.isCharging ? '⚡' : '🔋';
      if (state.battery.isCharging) {
        batteryBadge.classList.add('charging');
        batteryBadge.title = state.battery.percentage + '% (' + t.acPlugged + ')';
      } else {
        batteryBadge.classList.remove('charging');
        batteryBadge.title = state.battery.percentage + '% (' + t.batteryPower + ')';
      }
    } else {
      batPct.textContent = 'AC';
      batIcon.textContent = '🔌';
      batteryBadge.title = t.acPlugged;
    }
  }
}

// Toggle language event
btnLang.addEventListener('click', () => {
  currentLang = currentLang === 'vi' ? 'en' : 'vi';
  updateLanguageTexts();
  renderUI();
});

// Window controls: Minimize & Close to tray ("Show hidden icons")
if (btnWinMin) {
  btnWinMin.addEventListener('click', (e) => {
    e.stopPropagation();
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.MinimizeWindow) {
      window.go.main.App.MinimizeWindow();
    } else if (window.runtime && window.runtime.WindowHide) {
      window.runtime.WindowHide();
    }
  });
}

if (btnWinClose) {
  btnWinClose.addEventListener('click', (e) => {
    e.stopPropagation();
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.HideWindow) {
      window.go.main.App.HideWindow();
    } else if (window.runtime && window.runtime.WindowHide) {
      window.runtime.WindowHide();
    }
  });
}

// Interacting with Wails Go Backend
async function handleToggle() {
  const willBeAwake = !state.isAwake;

  if (willBeAwake) {
    playBatyTuongSound(); // ⚡ BẬT Ý TƯỞNG ACOUSTIC SPARK!
  } else {
    playSleepRestSound();  // 🌙 Calm resting tone
  }

  if (window.go && window.go.main && window.go.main.App) {
    try {
      if (state.isAwake) {
        state = await window.go.main.App.StopAwake();
      } else {
        const keepDisplay = state.mode !== 'system_only' && chkKeepDisplay.checked;
        const duration = state.mode === 'timed' ? selectedMinutes : 0;
        state = await window.go.main.App.StartAwake(
          state.mode,
          keepDisplay,
          duration,
          chkSafeguard.checked,
          5,
          chkAutoLock.checked
        );
      }
      renderUI();
    } catch (err) {
      console.error('Error toggling awake:', err);
    }
  } else {
    // Fallback simulation for dev/browser testing
    state.isAwake = !state.isAwake;
    renderUI();
  }
}

mainToggleBtn.addEventListener('click', handleToggle);

// Immediate Workstation Lock button
btnLockNow.addEventListener('click', async () => {
  if (window.go && window.go.main && window.go.main.App) {
    try {
      await window.go.main.App.LockNow();
    } catch (err) {
      console.error('Error locking workstation:', err);
    }
  }
});

// Smoothly update running session with new options without double audio restart
async function applyAwakeSettings() {
  if (!state.isAwake) return;
  if (window.go && window.go.main && window.go.main.App) {
    try {
      const keepDisplay = state.mode !== 'system_only' && chkKeepDisplay.checked;
      const duration = state.mode === 'timed' ? selectedMinutes : 0;
      state = await window.go.main.App.StartAwake(
        state.mode,
        keepDisplay,
        duration,
        chkSafeguard.checked,
        5,
        chkAutoLock.checked
      );
      renderUI();
    } catch (err) {
      console.error('Error applying awake settings:', err);
    }
  }
}

// Mode button clicks
Object.keys(modeButtons).forEach((m) => {
  modeButtons[m].addEventListener('click', () => {
    state.mode = m;
    if (m === 'system_only') {
      chkKeepDisplay.checked = false;
    } else {
      chkKeepDisplay.checked = true;
    }
    renderUI();

    // If currently awake, seamlessly re-apply with new mode
    if (state.isAwake) {
      applyAwakeSettings();
    }
  });
});

// Preset chips click
document.querySelectorAll('.preset-chip').forEach((chip) => {
  chip.addEventListener('click', () => {
    document.querySelectorAll('.preset-chip').forEach(c => c.classList.remove('active'));
    chip.classList.add('active');
    selectedMinutes = parseInt(chip.getAttribute('data-min'), 10);

    if (state.isAwake && state.mode === 'timed') {
      applyAwakeSettings();
    }
  });
});

// Safeguard & Display toggle changes
[chkKeepDisplay, chkSafeguard, chkAutoLock].forEach((chk) => {
  chk.addEventListener('change', () => {
    if (state.isAwake) {
      applyAwakeSettings();
    }
  });
});

// Listen to backend events if Wails runtime is loaded
if (window.runtime) {
  window.runtime.EventsOn('state_updated', (updatedState) => {
    state = updatedState;
    renderUI();
  });

  window.runtime.EventsOn('safeguard_triggered', (msg) => {
    toastNotice.textContent = `⚠️ ${msg}`;
    toastNotice.style.color = 'var(--accent-amber)';
  });

  window.runtime.EventsOn('timer_finished', (msg) => {
    toastNotice.textContent = `⏰ ${msg}`;
    toastNotice.style.color = 'var(--accent-emerald)';
  });
}

// Initial state fetch from Go backend
async function init() {
  updateLanguageTexts();
  if (window.go && window.go.main && window.go.main.App) {
    try {
      state = await window.go.main.App.GetState();
    } catch (e) {
      console.warn('Backend not ready yet:', e);
    }
  }
  renderUI();
}

init();
