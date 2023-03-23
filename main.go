package main

import (
	"embed"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var AccessToken string

func main() {
	// Create an instance of the app structure
	app := NewApp()
	AccessToken = uuid.NewString()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "XYHELPER",
		Width:  1050,
		Height: 649,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewAssetHandler(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		OnDomReady: app.DomReady,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
