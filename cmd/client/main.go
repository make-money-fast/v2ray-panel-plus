package main

import (
	"github.com/getlantern/systray"
	v2ray_panel_plus "v2ray-panel-plus"
	"v2ray-panel-plus/pkg/api/client"
	"v2ray-panel-plus/pkg/conf"
	"v2ray-panel-plus/pkg/menu"
	"v2ray-panel-plus/pkg/runtime"
)

func init() {
	conf.InitDefaultConfigFile()
	conf.InitLocalConfig()
}

func main() {
	systray.Run(func() {
		ico, _ := v2ray_panel_plus.StaticFS.ReadFile("static/favicon.png")
		go func() {
			menu.Init()
		}()
		systray.SetTemplateIcon(ico, ico)
		systray.SetTooltip("v2-client")
		runtime.AutoStart()
		client.StartHttpServer()
	}, nil)

}
