//go:build windows

package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004

	WM_TRAYICON      = 0x0400 + 100 // WM_USER + 100
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONUP     = 0x0205
	WM_DESTROY       = 0x0002

	MF_STRING    = 0x00000000
	MF_SEPARATOR = 0x00000800

	TPM_RETURNCMD = 0x0100
	TPM_NONOTIFY  = 0x0080

	IDI_APPLICATION = 32512

	CMD_SHOW   = 1001
	CMD_TOGGLE = 1002
	CMD_LOCK   = 1003
	CMD_QUIT   = 1004
)

var (
	shell32DLL          = syscall.NewLazyDLL("shell32.dll")
	procShellNotifyIcon = shell32DLL.NewProc("Shell_NotifyIconW")
	procExtractIcon     = shell32DLL.NewProc("ExtractIconW")

	user32DLL               = syscall.NewLazyDLL("user32.dll")
	procRegisterClassEx     = user32DLL.NewProc("RegisterClassExW")
	procCreateWindowEx      = user32DLL.NewProc("CreateWindowExW")
	procDefWindowProc       = user32DLL.NewProc("DefWindowProcW")
	procDestroyWindow       = user32DLL.NewProc("DestroyWindow")
	procGetMessage          = user32DLL.NewProc("GetMessageW")
	procTranslateMessage    = user32DLL.NewProc("TranslateMessage")
	procDispatchMessage     = user32DLL.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32DLL.NewProc("PostQuitMessage")
	procCreatePopupMenu     = user32DLL.NewProc("CreatePopupMenu")
	procAppendMenu          = user32DLL.NewProc("AppendMenuW")
	procTrackPopupMenuEx    = user32DLL.NewProc("TrackPopupMenuEx")
	procSetForegroundWindow = user32DLL.NewProc("SetForegroundWindow")
	procGetCursorPos        = user32DLL.NewProc("GetCursorPos")
	procDestroyMenu         = user32DLL.NewProc("DestroyMenu")
	procLoadIcon            = user32DLL.NewProc("LoadIconW")
	procLoadImage           = user32DLL.NewProc("LoadImageW")
	procFindWindow          = user32DLL.NewProc("FindWindowW")
	procSendMessage         = user32DLL.NewProc("SendMessageW")

	kernel32DLL           = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandle   = kernel32DLL.NewProc("GetModuleHandleW")
	procGetModuleFileName = kernel32DLL.NewProc("GetModuleFileNameW")
)

type NOTIFYICONDATAW struct {
	CbSize            uint32
	_                 uint32 // 64-bit alignment padding
	HWnd              uintptr
	UID               uint32
	UFlags            uint32
	UCallbackMessage  uint32
	_                 uint32 // 64-bit alignment padding
	HIcon             uintptr
	SzTip             [128]uint16
	DwState           uint32
	DwStateMask       uint32
	SzInfo            [256]uint16
	UTimeoutOrVersion uint32
	SzInfoTitle       [64]uint16
	DwInfoFlags       uint32
	GuidItem          [16]byte
	HBalloonIcon      uintptr
}

type POINT struct {
	X int32
	Y int32
}

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

var (
	trayApp           *App
	trayHwnd          uintptr
	trayHIcon         uintptr
	currentCupEdition = "gdu"
	trayMu            sync.Mutex
	trayCreated       bool
)

func initTray(a *App) {
	trayMu.Lock()
	defer trayMu.Unlock()

	if trayCreated {
		return
	}
	trayApp = a

	go runTrayMessageLoop()
}

func resolveAssetPath(rel string) string {
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	var exePathBuf [260]uint16
	procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&exePathBuf[0])), 260)
	exeDir := filepath.Dir(syscall.UTF16ToString(exePathBuf[:]))
	candidate := filepath.Join(exeDir, rel)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

func loadCupHIcon(hInstance uintptr, edition string) uintptr {
	const IMAGE_ICON = 1
	const LR_LOADFROMFILE = 0x00000010

	var candidates []string
	if edition == "eco" {
		candidates = []string{
			"CoffeeGDU icons\\Coffee-icon-leafCup.ico",
			"Coffee-icon-leafCup.ico",
			"CoffeeGDU icons\\Coffee-24h-GDU-64Gperpackage-icon.ico",
			"build\\windows\\icon.ico",
		}
	} else {
		candidates = []string{
			"CoffeeGDU icons\\Coffee-24h-GDU-64Gperpackage-icon.ico",
			"Coffee-24h-GDU-64Gperpackage-icon.ico",
			"build\\windows\\icon.ico",
			"CoffeeGDU icons\\Coffee-icon-leafCup.ico",
		}
	}

	for _, p := range candidates {
		resolved := resolveAssetPath(p)
		if resolved != "" {
			pathPtr, _ := syscall.UTF16PtrFromString(resolved)
			hIcon, _, _ := procLoadImage.Call(0, uintptr(unsafe.Pointer(pathPtr)), IMAGE_ICON, 16, 16, LR_LOADFROMFILE)
			if hIcon != 0 {
				return hIcon
			}
		}
	}

	// Try extracting embedded main icon from the running .exe
	var exePathBuf [260]uint16
	procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&exePathBuf[0])), 260)
	hIcon, _, _ := procExtractIcon.Call(hInstance, uintptr(unsafe.Pointer(&exePathBuf[0])), 0)
	if hIcon != 0 && hIcon != 1 {
		return hIcon
	}

	// Fallback to standard application icon
	hIcon, _, _ = procLoadIcon.Call(0, uintptr(IDI_APPLICATION))
	return hIcon
}

func loadTrayIcon(hInstance uintptr) uintptr {
	edition := "gdu"
	if trayApp != nil && trayApp.cliOpts.EcoCup {
		edition = "eco"
	}
	currentCupEdition = edition
	return loadCupHIcon(hInstance, edition)
}

func switchCupIcon(edition string) {
	trayMu.Lock()
	defer trayMu.Unlock()

	edition = strings.ToLower(edition)
	if edition != "eco" && edition != "gdu" {
		edition = "gdu"
	}
	currentCupEdition = edition

	hInstance, _, _ := procGetModuleHandle.Call(0)
	newHIcon := loadCupHIcon(hInstance, edition)
	if newHIcon != 0 {
		trayHIcon = newHIcon
	}

	// Update window taskbar/titlebar icon
	setWindowIcons("Coffee24h ☕ GDU")

	// Update System Tray icon if created
	if trayCreated && trayHwnd != 0 && trayHIcon != 0 {
		var nid NOTIFYICONDATAW
		nid.CbSize = uint32(unsafe.Sizeof(nid))
		nid.HWnd = trayHwnd
		nid.UID = 1
		nid.UFlags = NIF_ICON
		nid.HIcon = trayHIcon
		procShellNotifyIcon.Call(uintptr(NIM_MODIFY), uintptr(unsafe.Pointer(&nid)))
	}
}

func runTrayMessageLoop() {
	hInstance, _, _ := procGetModuleHandle.Call(0)
	className, _ := syscall.UTF16PtrFromString("Coffee24hTrayWindowClass")

	var wc WNDCLASSEXW
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(trayWndProc)
	wc.HInstance = hInstance
	wc.LpszClassName = className

	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	hwnd, _, _ := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(className)),
		0, 0, 0, 0, 0,
		0, 0, hInstance, 0,
	)
	if hwnd == 0 {
		return
	}

	trayHwnd = hwnd

	// Load icon: Try loading dedicated tray icon from CoffeeGDU icons, build folder, or exe
	trayHIcon = loadTrayIcon(hInstance)

	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = hwnd
	nid.UID = 1
	nid.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	nid.UCallbackMessage = WM_TRAYICON
	nid.HIcon = trayHIcon

	tip, _ := syscall.UTF16FromString("Coffee24h - Keep Awake")
	copy(nid.SzTip[:], tip)

	procShellNotifyIcon.Call(uintptr(NIM_ADD), uintptr(unsafe.Pointer(&nid)))
	trayCreated = true

	// Message loop
	var msg MSG
	for {
		r1, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if r1 == 0 || int32(r1) == -1 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func trayWndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_TRAYICON:
		switch lParam {
		case WM_LBUTTONUP, WM_LBUTTONDBLCLK:
			// Click or double click: Anchor to bottom-right corner and bring window to foreground
			if trayApp != nil && trayApp.ctx != nil {
				trayApp.domReady(trayApp.ctx)
				runtime.WindowShow(trayApp.ctx)
				runtime.WindowUnminimise(trayApp.ctx)
			}
		case WM_RBUTTONUP:
			// Right click: Show context menu
			showTrayMenu(hwnd)
		}
		return 0
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	r, _, _ := procDefWindowProc.Call(hwnd, uintptr(msg), wParam, lParam)
	return r
}

func showTrayMenu(hwnd uintptr) {
	if trayApp == nil || trayApp.ctx == nil {
		return
	}

	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(hwnd)

	hMenu, _, _ := procCreatePopupMenu.Call()
	if hMenu == 0 {
		return
	}
	defer procDestroyMenu.Call(hMenu)

	titleStr, _ := syscall.UTF16PtrFromString("☕ Coffee24h - Mở giao diện")
	procAppendMenu.Call(hMenu, MF_STRING, CMD_SHOW, uintptr(unsafe.Pointer(titleStr)))

	toggleStr, _ := syscall.UTF16PtrFromString("⚡ Bật / Tắt Awake (Toggle)")
	procAppendMenu.Call(hMenu, MF_STRING, CMD_TOGGLE, uintptr(unsafe.Pointer(toggleStr)))

	lockStr, _ := syscall.UTF16PtrFromString("🔒 Khóa máy ngay (Win + L)")
	procAppendMenu.Call(hMenu, MF_STRING, CMD_LOCK, uintptr(unsafe.Pointer(lockStr)))

	procAppendMenu.Call(hMenu, MF_SEPARATOR, 0, 0)

	quitStr, _ := syscall.UTF16PtrFromString("❌ Thoát ứng dụng (Exit)")
	procAppendMenu.Call(hMenu, MF_STRING, CMD_QUIT, uintptr(unsafe.Pointer(quitStr)))

	cmd, _, _ := procTrackPopupMenuEx.Call(
		hMenu,
		TPM_RETURNCMD|TPM_NONOTIFY,
		uintptr(pt.X),
		uintptr(pt.Y),
		hwnd,
		0,
	)

	switch cmd {
	case CMD_SHOW:
		trayApp.domReady(trayApp.ctx)
		runtime.WindowShow(trayApp.ctx)
		runtime.WindowUnminimise(trayApp.ctx)
	case CMD_TOGGLE:
		trayApp.ToggleAwake()
	case CMD_LOCK:
		_ = trayApp.LockNow()
	case CMD_QUIT:
		trayApp.Quit()
	}
}

func updateTrayStatus(isAwake bool) {
	trayMu.Lock()
	defer trayMu.Unlock()

	if !trayCreated || trayHwnd == 0 {
		return
	}

	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = trayHwnd
	nid.UID = 1
	nid.UFlags = NIF_TIP

	var tipText string
	prefix := "Coffee24h"
	if currentCupEdition == "eco" {
		prefix = "Coffee24h (Eco Student)"
	}
	if isAwake {
		tipText = prefix + ": Đang giữ máy thức (Awake)"
	} else {
		tipText = prefix + ": Cho phép nghỉ (Sleep Allowed)"
	}

	tip, _ := syscall.UTF16FromString(tipText)
	copy(nid.SzTip[:], tip)

	procShellNotifyIcon.Call(uintptr(NIM_MODIFY), uintptr(unsafe.Pointer(&nid)))
}

func removeTray() {
	trayMu.Lock()
	defer trayMu.Unlock()

	if !trayCreated || trayHwnd == 0 {
		return
	}

	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = trayHwnd
	nid.UID = 1

	procShellNotifyIcon.Call(uintptr(NIM_DELETE), uintptr(unsafe.Pointer(&nid)))
	procDestroyWindow.Call(trayHwnd)
	trayCreated = false
	trayHwnd = 0
}

func setWindowIcons(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return
	}
	const (
		WM_SETICON = 0x0080
		ICON_SMALL = 0
		ICON_BIG   = 1
	)
	if trayHIcon != 0 {
		procSendMessage.Call(hwnd, WM_SETICON, ICON_SMALL, trayHIcon)
		procSendMessage.Call(hwnd, WM_SETICON, ICON_BIG, trayHIcon)
	}
}

