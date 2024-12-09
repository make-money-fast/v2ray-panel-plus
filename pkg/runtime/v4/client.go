package v4

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

func (v4Server) Start(uri string) error {
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

func (v4Server) Stop() {
	mu.Lock()
	defer mu.Unlock()

	if _server != nil {
		_server.Close()
		_server = nil
	}

	closeStopChan()
}

func (s *v4Server) Reload(uri string) error {
	s.Stop()
	return s.Start(uri)
}

func (v4Server) IsRunning() bool {
	return _server != nil
}

type v4Server struct{}

func New() *v4Server {
	return &v4Server{}
}
