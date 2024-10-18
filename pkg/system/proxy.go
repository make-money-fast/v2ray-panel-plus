package system

import (
	"github.com/pkg/errors"
	"sync"
	"v2ray-panel-plus/pkg/conf"
)

var (
	mu          sync.RWMutex
	proxyStatus = Off
	proxyMode   int
)

const (
	Off = 0
	On  = 1
)

const (
	ModeNone   = 0
	ModeGlobal = 1
	ModePac    = 2
)

func GetProxyStatus() int {
	mu.RLock()
	defer mu.RUnlock()

	s := proxyStatus
	return s
}

func SetProxyStatus(status int) {
	mu.Lock()
	defer mu.Unlock()

	proxyStatus = status
}

func GetMode() int {
	mu.RLock()
	defer mu.RUnlock()

	s := proxyMode
	return s
}

func SetModeAndStatus(mode int, status int) {
	mu.Lock()
	defer mu.Unlock()

	proxyMode = mode
	proxyStatus = status
}

func SetNone() error {
	rs, err := conf.GetRunningStatus()
	if err != nil {
		return errors.Wrap(err, "读取配置失败")
	}
	err = UnSetProxy()
	if err != nil {
		return err
	}
	rs.ProxyStatus = Off
	rs.ProxyMode = ModeNone
	SetModeAndStatus(ModeNone, Off)
	return conf.SetRunningStatus(rs)
}

func SetGlobal() error {
	rs, err := conf.GetRunningStatus()
	if err != nil {
		return errors.Wrap(err, "读取配置失败")
	}
	local, err := conf.GetLocalConfig()
	if err != nil {
		return err
	}
	err = SetProxy(local.HttpProxy())
	if err != nil {
		return err
	}
	rs.ProxyMode = ModeGlobal
	rs.ProxyStatus = On
	SetModeAndStatus(ModeGlobal, On)
	return conf.SetRunningStatus(rs)
}

func SetPac(address string) error {
	rs, err := conf.GetRunningStatus()
	if err != nil {
		return errors.Wrap(err, "读取配置失败")
	}
	err = SetProxy(address)
	if err != nil {
		return err
	}
	rs.ProxyMode = ModePac
	rs.ProxyStatus = On
	SetModeAndStatus(ModePac, On)
	return conf.SetRunningStatus(rs)
}
