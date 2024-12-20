package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	v2ray_panel_plus "github.com/make-money-fast/v2ray-panel-plus"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/api/client"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/menu"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/pac"
	client2 "github.com/make-money-fast/v2ray-panel-plus/pkg/runtime/client"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/system"
	"time"
)

var (
	test bool
	ui   bool
)

func init() {
	flag.BoolVar(&test, "t", false, "testing mode")
	flag.BoolVar(&ui, "b", true, "show ui")
}

func main() {
	flag.Parse()

	if test {
		fmt.Println("testing ok")
		return
	}

	conf.InitDefaultConfigFile()
	conf.InitLocalConfig()
	pac.InitGfw(conf.GetGfwPath())
	conf.InitRunningStatus()

	if ui {
		run()
	} else {
		runNoUI()
	}
}

func runNoUI() {
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
}

func run() {
	systray.Run(func() {
		ico, _ := v2ray_panel_plus.StaticFS.ReadFile("static/favicon.png")
		go func() {
			menu.Init()
		}()
		go func() {
			runNoUI()
		}()
		systray.SetTemplateIcon(ico, ico)
		systray.SetTooltip("v2-client")
		client.StartHttpServer()
	}, nil)
}
