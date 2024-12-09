//go:build !linux

package menu

import (
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime"
	"github.com/pkg/browser"
	"os"
)

var (
	UIAddress       = ""
	ServerUIAddress = ""
	CloseChan       = make(chan struct{})
)

func Init() {
	panelCh := systray.AddMenuItem("客户端管理", "")
	serverCh := systray.AddMenuItem("服务端管理", "")
	exitCh := systray.AddMenuItem("退出", "")

	for {
		select {
		case <-panelCh.ClickedCh:
			browser.OpenURL(UIAddress)
		case <-exitCh.ClickedCh:
			shutdown()
		case <-serverCh.ClickedCh:
			browser.OpenURL(ServerUIAddress)
		case <-CloseChan:
			shutdown()
		}
	}
}

func shutdown() {
	runtime.Stop()
	systray.Quit()
	os.Exit(0)
}

func Shutdown() {
	close(CloseChan)
}
