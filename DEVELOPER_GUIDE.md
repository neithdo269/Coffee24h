# 📖 Coffee24h — Tài Liệu Kiến Trúc & Cẩm Nang Gỡ Lỗi Cho Lập Trình Viên
> **Dự án:** Coffee24h Keep-Awake Utility (Gia Dinh University - GDU Edition)  
> **Tác giả:** Neith Do & Antigravity  
> **Ngôn ngữ & Công nghệ:** Go 1.22+, Wails v2, Win32 API, HTML5 / CSS3 / Vanilla JS (Vite), Web Audio API

---

## 📑 Mục Lục
1. [Tổng Quan Kiến Trúc Dự Án](#1-tổng-quan-kiến-trúc-dự-án)
2. [Bài Học Gỡ Lỗi #1: Khử Răng Cưa Bo Tròn 4 Góc (Subpixel Anti-Aliasing vs 1-bit GDI)](#2-bài-học-gỡ-lỗi-1-khử-răng-cưa-bo-tròn-4-góc)
3. [Bài Học Gỡ Lỗi #2: Bố Cục Thích Ứng Cục Bộ (Scoped Adaptive Layout cho Hẹn Giờ)](#3-bài-học-gỡ-lỗi-2-bố-cục-thích-ứng-cục-bộ)
4. [Bài Học Gỡ Lỗi #3: Win32 Power API & Ghim Luồng Hệ Điều Hành (`runtime.LockOSThread`)](#4-bài-học-gỡ-lỗi-3-win32-power-api--ghim-luồng-hệ-điều-hành)
5. [Bài Học Gỡ Lỗi #4: Triệt Tiêu Race Condition Bằng Atomic `sync.Mutex`](#5-bài-học-gỡ-lỗi-4-triệt-tiêu-race-condition-bằng-atomic-syncmutex)
6. [Bài Học Gỡ Lỗi #5: Tự Phục Hồi Khay Hệ Thống Khi Explorer Crash (`TaskbarCreated`)](#6-bài-học-gỡ-lỗi-5-tự-phục-hồi-khay-hệ-thống-khi-explorer-crash)
7. [Bài Học Gỡ Lỗi #6: Kiến Trúc Clean Code Đổi Phiên Bản `GDU` ⇄ `GDU student`](#7-bài-học-gỡ-lỗi-6-kiến-trúc-clean-code-đổi-phiên-bản)
8. [Cẩm Nang Tái Sử Dụng Code Cho Các Dự Án Desktop Tương Lai](#8-cẩm-nang-tái-sử-dụng-code-cho-các-dự-án-desktop-tương-lai)

---

## 1. Tổng Quan Kiến Trúc Dự Án

Coffee24h được xây dựng theo mô hình **Hybrid Desktop Application**:
- **Backend (Go + Win32 API):** Điều khiển luồng thực thi OS, can thiệp trạng thái nguồn (`SetThreadExecutionState`), khóa màn hình (`LockWorkStation`), quản lý pin (`GetSystemPowerStatus`), và biểu tượng khay hệ thống (`Shell_NotifyIconW`).
- **Bridge (Wails v2 IPC):** Cầu nối hai chiều an toàn giữa Go và Frontend JavaScript thông qua Runtime Events và Bindings.
- **Frontend (Vanilla JS + Vite):** Giao diện kính mờ (Glassmorphism), tổng hợp âm thanh chuông hòa âm ngoại tuyến (Web Audio API), từ điển song ngữ động (i18n), thiết kế responsive 350x350px chuẩn 1:1.

```
┌────────────────────────────────────────────────────────┐
│                   Giao Diện Người Dùng                  │
│       Vite + Vanilla JS + CSS3 Glassmorphism UI        │
│          (Web Audio Synthesizer, i18n, DOM)            │
└───────────────────────────▲────────────────────────────┘
                            │ Wails v2 IPC (Events & Bind)
┌───────────────────────────▼────────────────────────────┐
│                  Go Application Core                   │
│           (App State Machine, Atomic Mutex)            │
├────────────────────────────────────────────────────────┤
│                 Tầng Windows Win32 API                 │
│  - SetThreadExecutionState (Chống ngủ máy)            │
│  - Shell_NotifyIconW (Khay hệ thống Show Hidden Icons) │
│  - LockWorkStation (Khóa máy bảo mật Win+L)            │
│  - DwmSetWindowAttribute (Bo góc phần cứng Win 11)     │
│  - MessageBeep (Chuông hệ thống Win32)                 │
└────────────────────────────────────────────────────────┘
```

---

## 2. Bài Học Gỡ Lỗi #1: Khử Răng Cưa Bo Tròn 4 Góc

### ⚠️ Hiện Tượng Lỗi
Khi làm tròn 4 góc của cửa sổ Frameless, người dùng phản ánh: *"Các góc pixel trông như bị mất điểm ảnh / vỡ hạt răng cưa"*.

### 🔍 Nguyên Nhân Kỹ Thuật Sâu Xa
- Trong Win32 cổ điển, các nhà phát triển thường dùng:
  ```go
  hrgn, _, _ := procCreateRoundRectRgn.Call(0, 0, width+1, height+1, 28, 28)
  procSetWindowRgn.Call(hwnd, hrgn, 1)
  ```
- **Hạn chế chết người của `SetWindowRgn`:** GDI Region là một **mặt nạ nhị phân 1-bit** (Pixel chỉ có giá trị 0 hoặc 1, hoàn toàn **không hỗ trợ kênh mờ Alpha / Anti-Aliasing**).
- Trong khi đó, trình duyệt Chromium WebView2 vẽ viền cong `border-radius: 14px` bằng kỹ thuật **Sub-pixel Anti-aliasing** (tạo các điểm ảnh bán trong suốt mờ dần 80%, 50%, 20% ở rìa cong).
- Mặt nạ 1-bit của `SetWindowRgn` đã chém đứt các điểm ảnh bán trong suốt này theo tọa độ số nguyên, biến đường cong mượt mà thành **các bậc thang răng cưa 1px**, tạo ảo giác điểm ảnh bị khuyết tật.

###  Giải Pháp Hiện Đại Chuẩn Công Nghiệp
1. **Xóa bỏ hoàn toàn mặt nạ 1-bit GDI:**
   ```go
   procSetWindowRgn.Call(hwnd, 0, 1) // Trả quyền dựng hình 100% cho WebView2
   ```
2. **Kích hoạt DirectComposition 32-bit ARGB trong Wails:**
   ```go
   BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0}, // Trong suốt tuyệt đối
   Windows: &windows.Options{
       WebviewIsTransparent:              true,
       WindowIsTranslucent:               true,
       DisableFramelessWindowDecorations: false,
   },
   ```
3. **Kích hoạt GPU Hardware Layer trong CSS:**
   ```css
   html {
     background: transparent !important;
     overflow: hidden;
   }
   body {
     border-radius: 14px;
     border: 1px solid rgba(59, 130, 246, 0.4);
     box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08); /* Bỏ bóng thừa ngoài rìa */
     -webkit-backface-visibility: hidden;
     backface-visibility: hidden;
     transform: translateZ(0); /* Thúc đẩy GPU render khử răng cưa mượt mà */
   }
   ```
4. **Windows 11 DWM Hardware Anti-aliasing:**
   ```go
   if procDwmSetWindowAttribute.Find() == nil {
       var preference uint32 = 2 // DWMWCP_ROUND
       _, _, _ = procDwmSetWindowAttribute.Call(hwnd, 33, uintptr(unsafe.Pointer(&preference)), 4)
   }
   ```

---

## 3. Bài Học Gỡ Lỗi #2: Bố Cục Thích Ứng Cục Bộ

### ⚠️ Hiện Tượng Lỗi
Khi bấm sang mục **Hẹn Giờ**, dòng chữ chân trang `"Tương thích hệ điều hành Windows Win32 API • Made with ❤️ for GDU"` bị biến mất.

### 🔍 Nguyên Nhân Kỹ Thuật
- Cửa sổ ứng dụng được cố định ở kích thước vuông chuẩn **350x350px** với `overflow: hidden`.
- Khi ở chế độ **Vô Hạn** hoặc **Chạy Nền**, thanh chọn thời gian nhanh (`presetsBar`) bị ẩn, tổng chiều cao nội dung là ~320px, dòng chân trang hiển thị trọn vẹn.
- Khi người dùng bấm sang **Hẹn Giờ**, thanh `presetsBar` hiển thị (`display: flex`) chiếm thêm ~24px chiều cao.
- Điều này đẩy thẻ cài đặt xuống dưới, khiến dòng chữ chân trang bị lọt ra ngoài mép 350px và bị cắt mất.
- **Tại sao không tăng chiều cao toàn bộ cửa sổ lên 370px?** Nếu tăng toàn cục, 2 chế độ "Vô Hạn" và "Chạy Nền" sẽ bị giãn khoảng trống ở đáy, làm mất đi tỷ lệ vuông vức 1:1 cân đối trên màn hình laptop 15.6" (như ASUS ROG GL552VX).

###  Giải Pháp Scoped Adaptive CSS
Giữ nguyên kích thước chuẩn **350x350px** và áp dụng quy tắc **Scoped CSS**: Chỉ kích hoạt thu gọn khoảng cách khi chế độ Hẹn Giờ được chọn, giữ nguyên 100% giao diện cho 2 chế độ còn lại.

```javascript
// main.js - Tự động gắn / gỡ class trạng thái
presetsBar.style.display = state.mode === 'timed' ? 'flex' : 'none';
if (state.mode === 'timed') {
  document.body.classList.add('mode-timed');
} else {
  document.body.classList.remove('mode-timed');
}
```

```css
/* style.css - Co gọn khoảng cách 23px cục bộ CHỈ KHI ở chế độ Hẹn Giờ */
body.mode-timed .app-container { gap: 3px; }
body.mode-timed .hero-toggle-card { padding: 3px 8px 4px 8px; }
body.mode-timed .main-toggle-btn { width: 36px; height: 36px; }
body.mode-timed .controls-card { padding: 4px 8px 3px 8px; gap: 3px; }
body.mode-timed .presets-bar { margin-block-start: 2px; }
body.mode-timed .preset-chip { padding-block: 2px; font-size: 9px; }
```

---

## 4. Bài Học Gỡ Lỗi #3: Win32 Power API & Ghim Luồng Hệ Điều Hành

### ⚠️ Hiểm Họa Ẩn Giấu
Hàm Win32 `SetThreadExecutionState` có đặc tính: **Nó gán cờ duy trì trạng thái thức trực tiếp vào luồng hệ điều hành (OS Thread) đang gọi nó**.
Nếu luồng đó kết thúc hoặc chuyển trạng thái, Windows sẽ **tự động thu hồi cờ Awake**, khiến máy tính đột ngột đi ngủ dù app vẫn đang chạy!

### 🔍 Đặc Tính Của Go Scheduler
Go sử dụng cơ chế lập lịch M:N (hàng nghìn goroutine chia sẻ một vài OS thread). Một goroutine đang chạy có thể bị Go runtime tạm dừng (park) hoặc di chuyển sang OS thread khác bất kỳ lúc nào. Nếu gọi `SetThreadExecutionState` trong một goroutine ngẫu nhiên, cờ Awake sẽ hoạt động chập chờn không thể đoán trước.

###  Giải Pháp: Pinning OS Thread
```go
// inhibitor_windows.go
func newPlatformInhibitor() *WindowsInhibitor {
    w := &WindowsInhibitor{
        reqChan: make(chan stateReq),
    }

    // Khởi chạy goroutine chuyên trách và GHIM CHẶT vào một OS thread duy nhất
    go func() {
        runtime.LockOSThread() // CỰC KỲ QUAN TRỌNG: Khóa chặt Goroutine vào OS Thread này
        defer runtime.UnlockOSThread()

        for req := range w.reqChan {
            r1, _, err := procSetThreadExecState.Call(req.flags)
            if r1 == 0 {
                req.resErr <- err
            } else {
                req.resErr <- nil
            }
        }
    }()

    return w
}
```

---

## 5. Bài Học Gỡ Lỗi #4: Triệt Tiêu Race Condition Bằng Atomic `sync.Mutex`

### ⚠️ Hiện Tượng Lỗi (Crash do TOCTOU)
Trong các ứng dụng đếm giờ, tại thời điểm **chính xác 00:00:00**, goroutine đếm giờ đang thực thi hàm `StopAwake`. Đúng phần nghìn giây đó, người dùng bấm vào nút bật/tắt (Toggle) trên giao diện.
Hai goroutine cùng ghi vào `a.state` và cùng thao tác đóng channel `a.timerDone` dẫn đến lỗi `panic: close of closed channel` hoặc trạng thái không đồng nhất.

###  Giải Pháp: Khóa Atomic Mutex Toàn Trình
```go
// app.go
func (a *App) ToggleAwake() (AwakeState, error) {
    a.mu.Lock()
    defer a.mu.Unlock()

    if a.state.IsAwake {
        return a.stopAwakeLocked()
    }
    return a.startAwakeLocked(a.state.Mode, a.state.KeepDisplay, a.state.TargetSeconds/60, a.state.BatteryGuard, a.state.SafeguardPct, a.state.AutoLockOnEnd)
}

func (a *App) stopTimerLocked() {
    if a.timerDone != nil {
        close(a.timerDone)
        a.timerDone = nil // Đặt lại nil an toàn để không bao giờ bị double-close
    }
}
```

---

## 6. Bài Học Gỡ Lỗi #5: Tự Phục Hồi Khay Hệ Thống Khi Explorer Crash

### ⚠️ Hiện Tượng Lỗi
Khi người dùng khởi động lại Windows Explorer (hoặc Explorer bị treo/crash rồi tự bật lại qua Task Manager), tất cả các biểu tượng khay hệ thống (System Tray) của ứng dụng bên thứ ba sẽ bị biến mất hoàn toàn.

###  Giải Pháp: Lắng Nghe Thông Điệp `TaskbarCreated`
Windows phát đi một thông điệp Broadcast Message tên là `TaskbarCreated` mỗi khi thanh Taskbar được tạo mới. Ta chỉ cần đăng ký lắng nghe thông điệp này:

```go
// tray_windows.go
func runTrayMessageLoop() {
    // Đăng ký nhận ID thông điệp TaskbarCreated từ hệ điều hành
    name, _ := syscall.UTF16PtrFromString("TaskbarCreated")
    r, _, _ := procRegisterWindowMessage.Call(uintptr(unsafe.Pointer(name)))
    wmTaskbarCreated = r
    // ...
}

func trayWndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
    // Khi Explorer sập và sống lại, tự động tái tạo icon vào khay
    if wmTaskbarCreated != 0 && uintptr(msg) == wmTaskbarCreated {
        var nid NOTIFYICONDATAW
        nid.CbSize = uint32(unsafe.Sizeof(nid))
        nid.HWnd = hwnd
        nid.UID = 1
        nid.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
        nid.UCallbackMessage = WM_TRAYICON
        nid.HIcon = trayHIcon
        procShellNotifyIcon.Call(uintptr(NIM_ADD), uintptr(unsafe.Pointer(&nid)))
        return 0
    }
    // ...
}
```

---

## 7. Bài Học Gỡ Lỗi #6: Kiến Trúc Clean Code Đổi Phiên Bản

### 🎯 Yêu Cầu Thiết Kế
Bấm vào chữ `GDU` (hoặc icon Ly cà phê) sẽ lập tức chuyển đổi mượt mà giữa:
- **Bản Giảng Viên:** Badge đỏ ngọn lửa `GDU`, Tách Xanh Logo GDU.
- **Bản Sinh Viên:** Badge xanh ngọc `GDU student`, Ly Tái Chế Thân Thiện Môi Trường.

###  Nguyên Tắc Thiết Kế Single Source of Truth
Thay vì viết các khối `if/else` rải rác trong nhiều sự kiện, toàn bộ cấu hình phiên bản được cô lập thành hằng số bất biến (`Object.freeze`):

```javascript
// main.js
const EDITIONS = Object.freeze({
  LECTURER: 'lecturer',
  STUDENT: 'student'
});

const EDITION_CONFIGS = Object.freeze({
  [EDITIONS.LECTURER]: {
    tagText: 'GDU',
    tagClass: 'gdu-tag',
    logo: coffeeLogo,
    cupBtnClass: '',
    titleVi: '☕ Phiên bản Giảng Viên - Bấm để đổi sang GDU student',
    toastVi: '☕ Đã chuyển sang GDU Giảng Viên (Tách Xanh Logo GDU)!'
  },
  [EDITIONS.STUDENT]: {
    tagText: 'GDU student',
    tagClass: 'eco-tag',
    logo: ecoLogo,
    cupBtnClass: 'eco',
    titleVi: '🌱 Phiên bản Sinh Viên - Bấm để đổi sang GDU Giảng Viên',
    toastVi: '🌱 Đã chuyển sang GDU Sinh Viên (Ly Tái Chế)!'
  }
});

function normalizeEdition(edition) {
  const norm = String(edition || '').toLowerCase().trim();
  return (norm === 'student' || norm === 'eco') ? EDITIONS.STUDENT : EDITIONS.LECTURER;
}
```

---

## 8. Cẩm Nang Tái Sử Dụng Code Cho Các Dự Án Desktop Tương Lai

Khi phát triển các ứng dụng Desktop Windows sử dụng **Go + Wails v2**:

| Thành Phần | Giải Pháp Tốt Nhất | Lý Do / Hiệu Quả |
| :--- | :--- | :--- |
| **Cửa sổ bo góc** | `BackgroundColour: {0,0,0,0}` + `WebviewIsTransparent: true` + CSS `transform: translateZ(0)` | Khử răng cưa bằng GPU 32-bit ARGB, không bao giờ dùng `SetWindowRgn` 1-bit |
| **Ghim luồng Win32** | `runtime.LockOSThread()` trong Goroutine riêng | Chống Go Scheduler di chuyển luồng làm mất cờ API hệ điều hành |
| **Khay hệ thống** | Đăng ký `RegisterWindowMessageW("TaskbarCreated")` | Tự động phục hồi icon nếu Explorer.exe bị khởi động lại |
| **An toàn luồng** | `sync.Mutex` cho mọi thao tác ghi State & đóng Channel | Triệt tiêu hoàn toàn lỗi panic và Race Condition |
| **Bố cục co giãn** | Dùng class CSS cục bộ (Scoped Class) như `.mode-timed` | Đảm bảo tính toán khoảng cách pixel hoàn hảo cho từng chế độ |
| **Âm thanh ngoại tuyến** | Web Audio API Oscillator (`triangle` / `sine`) | Không cần phụ thuộc file âm thanh `.wav`/`.mp3` cồng kềnh, chạy 100% offline |

---
*Tài liệu được khởi tạo và lưu trữ trong mã nguồn dự án Coffee24h.*
