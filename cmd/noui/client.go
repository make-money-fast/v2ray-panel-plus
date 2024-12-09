package main

import (
	"flag"
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/api/client"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/pac"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime"
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
	runtime.InitServer(runtime.V5)
	conf.InitDefaultConfigFile()
	conf.InitLocalConfig()
	pac.InitGfw(conf.GetGfwPath())
	conf.InitRunningStatus()
	runNoUI()
	client.StartHttpServer("linux", "linux")
}

func runNoUI() {
	// 检测上次运行配置.
	status, err := conf.GetRunningStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	if status.RunningUUID != "" {
		path, err := conf.ActiveRuntimeConfigFile(status.RunningUUID)
		if err != nil {
			return
		}
		err = runtime.Start(path)
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
