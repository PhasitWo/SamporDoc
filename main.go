package main

import (
	"embed"
	// "runtime"

	"github.com/wailsapp/wails/v2"
	// "github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	// wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "development"

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Menu
	// AppMenu := menu.NewMenu()
	// if runtime.GOOS == "darwin" {
	// 	AppMenu.Append(menu.AppMenu()) // On macOS platform, this must be done right after `NewMenu()`
	// }
	// FileMenu := AppMenu.AddSubmenu("Preferences")
	// FileMenu.AddText("Setting", nil, func(_ *menu.CallbackData) {
	// 	wailsRuntime.EventsEmit(app.ctx, "navigate", "/setting")
	// })

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "SamporDoc " + version,
		Width:            700,
		Height:           700,
		DisableResize:    true,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		// Menu:             AppMenu,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
