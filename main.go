package main

import (
	"embed"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	AccessToken string
	Version     = "dev"

	BaseURI = "http://freechat.lidong.xin"
	app     *App
)

func main() {
	var ctx = gctx.New()

	var title = "XYHELPER 开源免费的AI助理 https://xyhelper.cn version: " + Version
	latestVersion, err := GetLatestVersion(ctx)
	if err == nil {

		if gstr.CompareVersion(latestVersion, Version) > 0 {
			title = title + " 有新版本 " + latestVersion + " 请前往 https://xyhelper.cn 下载"
		}
	}

	// Create an instance of the app structure
	app = NewApp()
	AccessToken = uuid.NewString()
	// Create application with options
	err = wails.Run(&options.App{
		Title:  title,
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
	g.Log().Info(gctx.New(), "XYHELPER 开源免费的AI助理 https://xyhelper.cn version: "+Version)
}

func GetLatestVersion(ctx g.Ctx) (lasterVersion string, err error) {
	httpClient := g.Client()
	httpProxyAddr, err := g.Cfg().Get(ctx, "httpProxyAddr")
	if err == nil && httpProxyAddr.String() != "" {
		g.Log().Info(ctx, "httpProxyAddr", httpProxyAddr.String())
		httpClient.SetProxy(httpProxyAddr.String())
	}
	r, err := httpClient.Get(ctx, "https://xyhelper.cn/version.txt?time="+gconv.String(gtime.Timestamp()), nil)
	if err != nil {
		g.Log().Error(ctx, "GetLatestVersion", err.Error())
		return
	}
	if r.StatusCode != 200 {
		g.Log().Error(ctx, "GetLatestVersion", "StatusCode!=200")
		return "", gerror.New("StatusCode!=200")
	}
	lasterVersion = gstr.Trim(r.ReadAllString())
	return
}
