package main

import (
	"embed"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cliOpts := parseCLI()
	if cliOpts.ShowHelp {
		printHelp()
		return
	}
	if cliOpts.ShowVersion {
		printVersion()
		return
	}

	app := NewApp()
	app.cliOpts = cliOpts

	err := wails.Run(&options.App{
		Title:             "Coffee24h ☕ GDU",
		Width:             350,
		Height:            350,
		MinWidth:          320,
		MinHeight:         320,
		DisableResize:     false,
		Frameless:         true,
		StartHidden:       cliOpts.StartTray,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0}, // Transparent background
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: false,
		},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnShutdown:       app.shutdown,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "coffee24h-gdu-keepawake-2026",
			OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
				if app != nil && app.ctx != nil {
					for _, arg := range data.Args {
						lower := strings.ToLower(arg)
						if lower == "--student" || lower == "-e" || lower == "--eco" || lower == "--leaf" || lower == "--recycle" {
							app.SetCupEdition("student")
							break
						} else if lower == "--lecturer" || lower == "--teacher" || lower == "--gdu" {
							app.SetCupEdition("lecturer")
							break
						}
					}
					app.domReady(app.ctx)
					runtime.WindowShow(app.ctx)
					runtime.WindowUnminimise(app.ctx)
				}
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error starting Coffee24h:", err.Error())
	}
}
