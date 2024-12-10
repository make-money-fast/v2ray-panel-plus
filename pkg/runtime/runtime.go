package runtime

import (
	v5 "github.com/make-money-fast/v2ray-panel-plus/pkg/runtime/v5"
)

const (
	V4 = iota + 1
	V5
)

var (
	_server Server
)

type Server interface {
	Start(string) error
	Stop()
	Reload(string) error
	IsRunning() bool
}

func InitServer(version int) {
	switch version {
	//case V4:
	//_server = v4.New()
	default:
		_server = v5.New()
	}
}

func Start(uri string) error {
	return _server.Start(uri)
}

func Stop() {
	_server.Stop()
}

func Reload(uri string) error {
	return _server.Reload(uri)
}

func IsRunning() bool {
	return _server.IsRunning()
}
