package runtime

import (
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime/client"
	"log"
)

// AutoStart 自动开启，用于服务中断后的重启
func AutoStart() {
	active, err := conf.GetActiveConfig()
	if err != nil {
		log.Println("[warning] [自启]: GetActiveConfig ", err)
		return
	}
	path, err := conf.ActiveRuntimeConfigFile(active.UUID)
	if err != nil {
		log.Println("[warning] [自启]: ActiveRuntimeConfigFile ", err)
		return
	}
	_ = client.Start(path)
}
