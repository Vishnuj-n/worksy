package main

import (
	"embed"

	"focusplay/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	a := app.New()

	err := wails.Run(&options.App{
		Title:            "FocusPlay",
		Width:            480,
		Height:           660,
		MinWidth:         420,
		MinHeight:        580,
		Frameless:        false,
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 26, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: a.Startup,
		Bind: []interface{}{
			a,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
