package client

import (
	_ "embed"
	core "github.com/clearcodecn/v2ray-core"
	"github.com/clearcodecn/v2ray-core/v2raystart"
	"sync"
)

var (
	mu       sync.Mutex
	stopChan chan struct{}
	_server  core.Server
	Version  = "v4.23.2"
)

func closeStopChan() {
	select {
	case <-stopChan:
	default:
		if stopChan != nil {
			close(stopChan)
			stopChan = nil
		}
	}
}

func Start(uri string) error {
	mu.Lock()
	defer mu.Unlock()
	closeStopChan()

	stopChan = make(chan struct{})

	server, err := v2raystart.Start(uri, stopChan)
	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	_server = server
	return nil
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if _server != nil {
		_server.Close()
		_server = nil
	}

	closeStopChan()
}

func Reload(uri string) error {
	Stop()
	return Start(uri)
}

func IsRunning() bool {
	return _server != nil
}
