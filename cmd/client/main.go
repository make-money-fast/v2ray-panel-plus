package main

import (
	"github.com/getlantern/systray"
	"time"
	v2ray_panel_plus "v2ray-panel-plus"
	"v2ray-panel-plus/pkg/api/client"
	"v2ray-panel-plus/pkg/conf"
	"v2ray-panel-plus/pkg/menu"
	"v2ray-panel-plus/pkg/pac"
	client2 "v2ray-panel-plus/pkg/runtime/client"
	"v2ray-panel-plus/pkg/system"
)

func init() {
	conf.InitDefaultConfigFile()
	conf.InitLocalConfig()
	pac.InitGfw(conf.GetGfwPath())
	conf.InitRunningStatus()
}

func main() {
	run()
}

func run() {
	systray.Run(func() {
		ico, _ := v2ray_panel_plus.StaticFS.ReadFile("static/favicon.png")
		go func() {
			menu.Init()
		}()
		go func() {
			time.Sleep(2 * time.Second)
			// 检测上次运行配置.
			status, err := conf.GetRunningStatus()
			if err != nil {
				return
			}
			if status.RunningUUID != "" {
				path, err := conf.ActiveRuntimeConfigFile(status.RunningUUID)
				if err != nil {
					return
				}
				err = client2.Start(path)
				if err != nil {
					return
				}
			}
			if status.ProxyStatus == system.Off {
				return
			}
			switch status.ProxyMode {
			case system.ModeNone:
				err = system.SetNone()
			case system.ModePac:
				for {
					addr := client.GetPacAddress()
					if addr == "" {
						time.Sleep(1 * time.Second)
						continue
					}
					err = system.SetPac(addr)
					break
				}
			case system.ModeGlobal:
				err = system.SetGlobal()
			}
		}()
		systray.SetTemplateIcon(ico, ico)
		systray.SetTooltip("v2-client")
		client.StartHttpServer()
	}, nil)
}
