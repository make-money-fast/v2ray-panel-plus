package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	myIP string
)

func GetMyIP() string {
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
