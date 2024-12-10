package helpers

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

func CheckPort(addr interface{}) bool {
	var dst string
	switch v := addr.(type) {
	case string:
		dst = v
	case int:
		dst = fmt.Sprintf("127.0.0.1:%d", v)
	}
	conn, err := net.Dial("tcp", dst)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func CheckPortbyProtocol(protocol string, addr interface{}) bool {
	var dst string
	switch v := addr.(type) {
	case string:
		dst = v
	case int:
		dst = fmt.Sprintf("127.0.0.1:%d", v)
	}
	conn, err := net.Dial(protocol, dst)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func CheckPorxy(u string, dst string) bool {
	proxyUrl, _ := url.Parse(u)

	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: 5 * time.Second,
	}

	resp, err := cli.Get(dst)
	if err != nil {
		return false
	}
	resp.Body.Close()

	return true
}
