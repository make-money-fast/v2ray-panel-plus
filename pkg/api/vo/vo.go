package vo

import (
	"encoding/base64"
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"strconv"
)

var (
	Token = base64.StdEncoding.EncodeToString([]byte("helloworld"))
)

type RequestContext struct {
	ServerUrl string `json:"server_url"`
}

func (c *RequestContext) Url(path string) string {
	return c.ServerUrl + "/api" + path
}

type ClientEditConfigRequest struct {
	RequestContext
	EditConfigRequest
}

type EditConfigRequest struct {
	UUID     string `json:"uuid"`     // 本地配置的唯一id
	Alias    string `json:"alias"`    // 别名
	Vmess    string `json:"vmess"`    // vmess
	Host     string `json:"host"`     // 地址
	Port     string `json:"port"`     // 端口
	Protocol string `json:"protocol"` // 协议.
	Id       string `json:"id"`       // 用户id
	Status   int    `json:"status"`   // 状态
	Ts       int64  `json:"ts"`       // 创建的时间戳
	Type     string `json:"type"`
}

func (req *EditConfigRequest) ToConfig() *conf.ServerConfig {
	item := &conf.ServerConfig{
		UUID: req.UUID,
	}
	port, _ := strconv.Atoi(req.Port)
	item.Alias = req.Alias
	item.Port = req.Port
	item.Protocol = req.Protocol
	item.Id = req.Id
	item.Host = conf.GetIP()
	item.Type = req.Type
	item.Config = &conf.ServerInbound{
		Tag:      fmt.Sprintf("%s_%d", req.Protocol, port),
		Port:     port,
		Protocol: "vmess",
		Settings: conf.Settings{
			Clients: []conf.Client{
				{
					Id:      req.Id,
					Level:   1,
					AlterId: 0,
				},
			},
		},
		StreamSettings: nil,
		Sniffing: conf.Sniffing{
			Enabled: true,
			DestOverride: []string{
				"tls",
				"http",
			},
		},
	}
	stream := &conf.StreamSetting{
		Network: item.Protocol,
	}
	if item.Protocol == "tcp" {
		stream.Network = "tcp"
		stream.TCPConfig = &conf.TCPConfig{
			HeaderConfig: map[string]string{
				"type": "none",
			},
		}
	}
	if item.Protocol == "ws" {
		stream.WSConfig = &conf.WebSocketConfig{}
	}
	/**
	"mtu": 1350,
	          "tti": 20,
	          "uplinkCapacity": 5,
	          "downlinkCapacity": 20,
	          "congestion": false,
	          "readBufferSize": 1,
	          "writeBufferSize": 1,
	          "header": {
	            "type": "wechat-video"
	          },
	          "seed": ""
	*/
	if item.Protocol == "kcp" || item.Protocol == "mkcp" {
		stream.KCPConfig = &conf.KCPConfig{
			Mtu:             1350,
			Tti:             20,
			UpCap:           10,
			DownCap:         100,
			Congestion:      true,
			ReadBufferSize:  10,
			WriteBufferSize: 10,
			HeaderConfig: map[string]string{
				"type": req.Type,
			},
		}
	}
	item.Config.StreamSettings = stream

	item.Vmess = item.BuildVmess()
	return item
}

type ClientDelConfigRequest struct {
	RequestContext
	DelConfigRequest
}

type DelConfigRequest struct {
	UUID string `json:"uuid"`
}
