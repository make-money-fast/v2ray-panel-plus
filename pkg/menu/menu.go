package menu

import (
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime/client"
	"github.com/pkg/browser"
	"os"
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
	client.Stop()
	systray.Quit()
	os.Exit(0)
}

func Shutdown() {
	close(CloseChan)
}
