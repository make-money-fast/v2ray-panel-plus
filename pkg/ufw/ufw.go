package ufw

import (
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/samber/lo"
	"log"
	"os"
	"os/exec"
)

func ActiveUFW() {
	serverConfig, err := conf.MergeServerConfig()
	if err != nil {
		log.Println("mergeServerConfig err:", err)
		return
	}
	var ports []int
	lo.ForEach(serverConfig.Inbounds, func(item conf.ServerInbound, index int) {
		ports = append(ports, item.Port)
	})
	lo.ForEach(ports, func(item int, index int) {
		openufw(item)
	})
}

func openufw(port int) {
	cmd := exec.Command("ufw", "allow", fmt.Sprintf("%d", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println("open ufw failed", err)
	}
}
