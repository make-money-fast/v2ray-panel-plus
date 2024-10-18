package conf

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"v2ray-panel-plus/pkg/helpers"
)

//go:embed localconfig.json
var defaultLocalConfig string

type LocalConfig struct {
	SocksAddress string `json:"socksAddress"`
	SocksPort    int    `json:"socksPort"`
	HttpAddress  string `json:"httpAddress"`
	HttpPort     int    `json:"httpPort"`
}

func (c LocalConfig) HttpProxy() string {
	return fmt.Sprintf("http://%s:%d", c.HttpAddress, c.HttpPort)
}

func (c LocalConfig) SocksPorxy() string {
	return fmt.Sprintf("socks5://%s:%d", c.SocksAddress, c.SocksPort)
}

func InitLocalConfig() {
	var local LocalConfig
	if err := helpers.ReadJSONFile(getLocalConfigPath(), &local); err != nil {
		json.Unmarshal([]byte(defaultLocalConfig), &local)
		if err := helpers.WriteJSONFile(getLocalConfigPath(), local); err != nil {
			panic(err)
		}
	}
}

func GetLocalConfig() (*LocalConfig, error) {
	var local LocalConfig
	if err := helpers.ReadJSONFile(getLocalConfigPath(), &local); err != nil {
		return nil, err
	}
	return &local, nil
}

func SetLocalConfig(config *LocalConfig) error {
	return helpers.WriteJSONFile(getLocalConfigPath(), config)
}
