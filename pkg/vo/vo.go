package vo

type EditConfigRequest struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Id       string `json:"id"`
	UUID     string `json:"uuid"`
	Network  string `json:"network"`
	Alias    string `json:"alias"`
	Settings struct {
		TcpSettings struct {
		} `json:"tcpSettings"`
		KcpSettings struct {
			UplinkCapacity   int  `json:"uplinkCapacity"`
			DownlinkCapacity int  `json:"downlinkCapacity"`
			Congestion       bool `json:"congestion"`
			Header           struct {
				Type string `json:"type"`
			} `json:"header"`
		} `json:"kcpSettings"`
		WsSettings struct {
			Path string `json:"path"`
		} `json:"wsSettings"`
	} `json:"settings"`
}

type UUIDConfigRequest struct {
	UUID string `json:"uuid"`
}

type ClientState struct {
	IsRunning       bool `json:"isRunning"`
	Socks           bool `json:"socks"`
	Http            bool `json:"http"`
	ConnectToServer bool `json:"connectToServer"`
	PorxyOK         bool `json:"porxyOK"`
}
