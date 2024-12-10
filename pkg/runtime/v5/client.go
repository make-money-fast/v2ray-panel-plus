package v5

import (
	core "github.com/make-money-fast/v2ray-core-v5"
	"github.com/make-money-fast/v2ray-core-v5/v2start"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	_server core.Server
)

func (v5Server) Start(uri string) error {
	mu.Lock()
	defer mu.Unlock()

	server, err := v2start.Start(uri)
	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	_server = server
	return nil
}

func (v5Server) Stop() {
	mu.Lock()
	defer mu.Unlock()

	if _server != nil {
		_server.Close()
		_server = nil
	}
}

func (s *v5Server) Reload(uri string) error {
	s.Stop()
	time.Sleep(300 * time.Millisecond)
	return s.Start(uri)
}

func (v5Server) IsRunning() bool {
	return _server != nil
}

type v5Server struct {
}

func New() *v5Server {
	return &v5Server{}
}
