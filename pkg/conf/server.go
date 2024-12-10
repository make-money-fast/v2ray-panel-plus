package conf

import (
	"encoding/json"
	"fmt"
	"github.com/clearcodecn/vmess"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/helpers"
	"github.com/samber/lo"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type ServerInbound struct {
	Tag            string         `json:"tag"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings *StreamSetting `json:"streamSettings,omitempty"`
	Sniffing       `json:"sniffing"`
}

type Client struct {
	Id      string `json:"id"`
	Level   int    `json:"level,omitempty"`
	AlterId int    `json:"alterId,omitempty"`
}

type Settings struct {
	Clients []Client `json:"clients"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

type ServerOutbound struct {
	Protocol string `json:"protocol"`
	Settings struct {
		DomainStrategy string `json:"domainStrategy,omitempty"`
	} `json:"settings"`
	Tag string `json:"tag"`
}

type ServerV2rayConfig struct {
	Log struct {
		Access   string `json:"access"`
		Error    string `json:"error"`
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds  []ServerInbound  `json:"inbounds"`
	Outbounds []ServerOutbound `json:"outbounds"`
	Dns       struct {
		Servers []string `json:"servers"`
	} `json:"dns"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			Ip          []string `json:"ip,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			InboundTag  []string `json:"inboundTag,omitempty"`
			Domain      []string `json:"domain,omitempty"`
			Protocol    []string `json:"protocol,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	Transport struct {
		KcpSettings struct {
			UplinkCapacity   int  `json:"uplinkCapacity"`
			DownlinkCapacity int  `json:"downlinkCapacity"`
			Congestion       bool `json:"congestion"`
		} `json:"kcpSettings"`
	} `json:"transport"`
}

type ServerConfig struct {
	UUID     string         `json:"uuid"`     // 本地配置的唯一id
	Config   *ServerInbound `json:"config"`   // json配置.
	Alias    string         `json:"alias"`    // 别名
	Vmess    string         `json:"vmess"`    // vmess
	Host     string         `json:"host"`     // 地址
	Port     string         `json:"port"`     // 端口
	Protocol string         `json:"protocol"` // 协议.
	Id       string         `json:"id"`       // 用户id
	Ts       int64          `json:"ts"`       // 创建的时间戳
	Type     string         `json:"type"`
}

func (c *ServerConfig) BuildVmess() string {
	ip := GetIP()

	link := vmess.Link{
		Version: 0,
		Name:    fmt.Sprintf("%s_%s", ip, c.Port),
		Address: ip,
		Port:    c.Port,
		Id:      c.Id,
		Aid:     "0",
		Network: c.Protocol,
		Type:    "",
	}

	if c.Config.StreamSettings.KCPConfig != nil {
		link.Type = c.Config.StreamSettings.KCPConfig.HeaderConfig["type"]
	}

	return link.ToVmessLink()
}

func GetIP() string {
	var data = make(map[string]string)
	rsp, err := http.Get("http://ipinfo.io")
	if err != nil {
		return ""
	}
	jsondata, _ := ioutil.ReadAll(rsp.Body)
	json.Unmarshal(jsondata, &data)
	return data["ip"]
}

func MergeServerConfig() (*ServerV2rayConfig, error) {
	serverConfig, err := GetServerConfigList()
	if err != nil {
		return nil, err
	}
	var inBounds []ServerInbound
	lo.ForEach(serverConfig, func(item *ServerConfig, index int) {
		inBounds = append(inBounds, *item.Config)
	})
	var cfg ServerV2rayConfig
	template := defaultServerTemplate()
	err = helpers.ReadJSONFile(template, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.Log.Error = filepath.Join(defaultConfigDirectory(), "error.log")
	cfg.Log.Access = filepath.Join(defaultConfigDirectory(), "access.log")
	cfg.Inbounds = inBounds
	return &cfg, nil
}
