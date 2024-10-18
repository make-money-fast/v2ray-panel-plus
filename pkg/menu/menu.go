package menu

import (
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"os"
	"v2ray-panel-plus/pkg/conf"
	"v2ray-panel-plus/pkg/runtime/client"
)

var (
	UIAddress = ""
	CloseChan = make(chan struct{})
)

func Init() {
	panelCh := systray.AddMenuItem("管理后台", "")
	exitCh := systray.AddMenuItem("退出", "")

	for {
		select {
		case <-panelCh.ClickedCh:
			browser.OpenURL(UIAddress)
		case <-exitCh.ClickedCh:
			shutdown()
		case <-CloseChan:
			shutdown()
		}
	}
}

func shutdown() {
	list, _ := conf.GetConfigList()
	for _, item := range list {
		item.Status = conf.StatusDown
	}
	conf.SaveConfigList(list)
	client.Stop()

	systray.Quit()
	os.Exit(0)
}

func Shutdown() {
	close(CloseChan)
}
