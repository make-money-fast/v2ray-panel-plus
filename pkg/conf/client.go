package conf

import (
	_ "embed"
	"encoding/json"
	"github.com/clearcodecn/vmess"
	"path/filepath"
	"strconv"
)

//go:embed default_client_config.json
var ClientDefaultConfig string

type ClientInBound struct {
	Tag      string `json:"tag"`
	Port     int    `json:"port"`
	Listen   string `json:"listen"`
	Protocol string `json:"protocol"`
	Sniffing struct {
		Enabled      bool     `json:"enabled"`
		DestOverride []string `json:"destOverride"`
	} `json:"sniffing"`
	Settings struct {
		Auth             string `json:"auth"`
		Udp              bool   `json:"udp"`
		AllowTransparent bool   `json:"allowTransparent"`
	} `json:"settings"`
}

type ClientConfigV2ray struct {
	Log struct {
		Access   string `json:"access"`
		Error    string `json:"error"`
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds  []ClientInBound `json:"inbounds"`
	Outbounds []struct {
		Tag      string `json:"tag"`
		Protocol string `json:"protocol"`
		Settings struct {
			Vnext []struct {
				Address string `json:"address"`
				Port    int    `json:"port"`
				Users   []struct {
					Id       string `json:"id"`
					AlterId  int    `json:"alterId"`
					Email    string `json:"email"`
					Security string `json:"security"`
				} `json:"users"`
			} `json:"vnext,omitempty"`
			Response struct {
				Type string `json:"type"`
			} `json:"response,omitempty"`
		} `json:"settings"`
		StreamSettings *StreamSetting `json:"streamSettings,omitempty"`
		Mux            *struct {
			Enabled     bool `json:"enabled,omitempty"`
			Concurrency int  `json:"concurrency,omitempty"`
		} `json:"mux,omitempty"`
	} `json:"outbounds"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			InboundTag  []string `json:"inboundTag,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			Enabled     bool     `json:"enabled"`
			Domain      []string `json:"domain,omitempty"`
			Ip          []string `json:"ip,omitempty"`
			Port        string   `json:"port,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	UUID string `json:"uuid"` // 本地记录的uuid
}

type StreamSetting struct {
	Network    string              `json:"network"`
	Security   string              `json:"security"`
	TCPConfig  *TCPConfig          `json:"tcpSettings"`
	KCPConfig  *KCPConfig          `json:"kcpSettings"`
	WSConfig   *WebSocketConfig    `json:"wsSettings"`
	HTTPConfig *HTTPConfig         `json:"httpSettings"`
	DSConfig   *DomainSocketConfig `json:"dsSettings"`
	QUICConfig *QUICConfig         `json:"quicSettings"`
}

type TCPConfig struct {
	HeaderConfig map[string]string `json:"header"`
}

type KCPConfig struct {
	Mtu             *uint32           `json:"mtu"`
	Tti             *uint32           `json:"tti"`
	UpCap           *uint32           `json:"uplinkCapacity"`
	DownCap         *uint32           `json:"downlinkCapacity"`
	Congestion      *bool             `json:"congestion"`
	ReadBufferSize  *uint32           `json:"readBufferSize"`
	WriteBufferSize *uint32           `json:"writeBufferSize"`
	HeaderConfig    map[string]string `json:"header"`
}

type WebSocketConfig struct {
	Path    string            `json:"path"`
	Path2   string            `json:"Path"` // The key was misspelled. For backward compatibility, we have to keep track the old key.
	Headers map[string]string `json:"headers"`
}

type HTTPConfig struct {
	Host *StringList `json:"host"`
	Path string      `json:"path"`
}

type StringList []string

type DomainSocketConfig struct {
	Path     string `json:"path"`
	Abstract bool   `json:"abstract"`
}

type QUICConfig struct {
	Header   map[string]string `json:"header"`
	Security string            `json:"security"`
	Key      string            `json:"key"`
}

const (
	HeaderNone        string = "none"         // 默认值，不进行伪装，发送的数据是没有特征的数据包。
	HeaderSrtp               = "srtp"         // 伪装成 SRTP 数据包，会被识别为视频通话数据（如 FaceTime）。
	HeaderUTP                = "utp"          // 伪装成 uTP 数据包，会被识别为 BT 下载数据。
	HeaderWechatVideo        = "wechat-video" // 伪装成微信视频通话的数据包。
	HeaderDTLS               = "dtls"         // 伪装成 DTLS 1.2 数据包。
	HeaderWireguard          = "wireguard"    // 伪装成 WireGuard 数据包。（并不是真正的 WireGuard 协议）
)

type HeaderSettings struct {
	Type string `json:"type"`
}

const (
	StatusStart = 1
	StatusDown  = 0
)

type ClientConfig struct {
	UUID     string             `json:"uuid"`     // 本地配置的唯一id
	Config   *ClientConfigV2ray `json:"config"`   // json配置.
	Alias    string             `json:"alias"`    // 别名
	Vmess    string             `json:"vmess"`    // vmess
	Host     string             `json:"host"`     // 地址
	Port     string             `json:"port"`     // 端口
	Protocol string             `json:"protocol"` // 协议.
	Id       string             `json:"id"`       // 用户id
	Status   int                `json:"status"`   // 状态
	Ts       int64              `json:"ts"`       // 创建的时间戳
}

func (c *ClientConfig) ReadLocalConfig() error {
	localConfig, err := GetLocalConfig()
	if err != nil {
		return err
	}
	inBounds := c.Config.Inbounds
	for idx, inbound := range inBounds {
		if inbound.Protocol == "socks" {
			inbound.Port = localConfig.SocksPort
			inbound.Listen = localConfig.SocksAddress
		}
		if inbound.Protocol == "http" {
			inbound.Port = localConfig.HttpPort
			inbound.Listen = localConfig.HttpAddress
		}
		inBounds[idx] = inbound
	}
	c.Config.Inbounds = inBounds
	c.Config.Log.Access = filepath.Join(defaultConfigDirectory(), "access.log")
	c.Config.Log.Error = filepath.Join(defaultConfigDirectory(), "error.log")
	return nil
}

func (c *ClientConfig) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *ClientConfigV2ray) GetVmess(name string) string {
	if len(c.Outbounds) == 0 {
		return ""
	}
	if len(c.Outbounds[0].Settings.Vnext) == 0 {
		return ""
	}
	if len(c.Outbounds[0].Settings.Vnext[0].Users) == 0 {
		return ""
	}
	if c.Outbounds[0].StreamSettings == nil {
		return ""
	}
	var link vmess.Link
	ip := c.Outbounds[0].Settings.Vnext[0].Address
	port := c.Outbounds[0].Settings.Vnext[0].Port
	uid := c.Outbounds[0].Settings.Vnext[0].Users[0].Id
	net := c.Outbounds[0].StreamSettings.Network

	link.Name = name
	link.Address = ip
	link.Port = strconv.Itoa(port)
	link.Id = uid
	link.Network = net
	link.Aid = strconv.Itoa(c.Outbounds[0].Settings.Vnext[0].Users[0].AlterId)
	if c.Outbounds[0].StreamSettings.KCPConfig != nil {
		link.Type = c.Outbounds[0].StreamSettings.KCPConfig.HeaderConfig["type"]
	}
	if c.Outbounds[0].StreamSettings.WSConfig != nil {
		link.Path = c.Outbounds[0].StreamSettings.WSConfig.Path
	}

	return link.ToVmessLink()
}

func (c *ClientConfigV2ray) GetServerAddressAndPort() (string, string) {
	if len(c.Outbounds) > 0 && len(c.Outbounds[0].Settings.Vnext) > 0 {
		addr := c.Outbounds[0].Settings.Vnext[0].Address
		pot := c.Outbounds[0].Settings.Vnext[0].Port
		return addr, strconv.Itoa(pot)
	}
	return "", ""
}

func (c *ClientConfigV2ray) GetProtocol() string {
	if len(c.Outbounds) > 0 {
		return c.Outbounds[0].StreamSettings.Network
	}
	return ""
}

func (c *ClientConfigV2ray) GetId() string {
	if len(c.Outbounds) > 0 && len(c.Outbounds[0].Settings.Vnext) > 0 {
		return c.Outbounds[0].Settings.Vnext[0].Users[0].Id
	}
	return ""
}
func ClientConfigFromString(str string) (*ClientConfig, error) {
	var conf ClientConfig
	err := json.Unmarshal([]byte(str), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func ClientConfigV2rayFromString(str string) (*ClientConfigV2ray, error) {
	var conf ClientConfigV2ray
	err := json.Unmarshal([]byte(str), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
