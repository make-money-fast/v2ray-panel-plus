package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var (
	myIP string
)

// GetPublicIP 获取公网ip
func GetPublicIP() string {
	if myIP != "" {
		return myIP
	}
	resp, _ := http.Get("https://ipinfo.io/json")
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var ret = make(map[string]string)
	json.Unmarshal(data, &ret)
	if ret["ip"] != "" {
		myIP = ret["ip"]
	}
	return myIP
}

// GetInternalIP 获取内网ip
func GetInternalIP() []string {
	faceAddress, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	var ipStrings []string

	for _, add := range faceAddress {
		ip, _, err := net.ParseCIDR(add.String())
		if err != nil {
			continue
		}
		if ip.IsLoopback() {
			continue
		}
		ipString := ip.String()
		if strings.Contains(ipString, ":") {
			continue
		}
		ipStrings = append(ipStrings, ipString)
	}
	return ipStrings
}
