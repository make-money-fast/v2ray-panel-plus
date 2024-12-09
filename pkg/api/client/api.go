package client

import (
	"encoding/json"
	"fmt"
	"github.com/clearcodecn/vmess"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/helpers"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/menu"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/system"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/vo"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/skip2/go-qrcode"
	"path/filepath"
	"strconv"
	"time"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func sendError(ctx *gin.Context, err error) {
	ctx.JSON(200, &Response[any]{
		Code: -1,
		Msg:  err.Error(),
	})
}

func sendSuccess(ctx *gin.Context, data any) {
	ctx.JSON(200, &Response[any]{
		Code: 0,
		Data: data,
	})
}

func initDefaultConfig(ctx *gin.Context) {
	newConf, err := conf.ClientConfigV2rayFromString(conf.ClientDefaultConfig)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "初始化配置失败"))
		return
	}
	newConf.UUID = uuid.New().String()
	addr, port := newConf.GetServerAddressAndPort()
	name := "初始化-" + helpers.NowDateTime()
	clientConfig := &conf.ClientConfig{
		UUID:     newConf.UUID,
		Config:   newConf,
		Alias:    name,
		Vmess:    newConf.GetVmess(name),
		Host:     addr,
		Port:     port,
		Protocol: newConf.GetProtocol(),
		Id:       newConf.GetId(),
		Status:   0,
		Ts:       time.Now().Unix(),
	}

	if err := conf.CreateOneConfig(clientConfig); err != nil {
		sendError(ctx, errors.Wrap(err, "保存配置失败"))
		return
	}
	sendSuccess(ctx, true)
}

func getConfigList(ctx *gin.Context) {
	list, _ := conf.GetConfigList()
	sendSuccess(ctx, list)
}

func editConfig(ctx *gin.Context) {
	var req vo.EditConfigRequest
	if err := ctx.BindJSON(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}
	config, err := conf.GetConfigByUUID(req.UUID)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取配置失败"))
		return
	}
	config.Host = req.Host
	config.Port = req.Port
	config.Id = req.Id
	config.Protocol = req.Network
	config.Alias = req.Alias

	config.Config.Outbounds[0].Settings.Vnext[0].Port = helpers.Str2Int(config.Port)
	config.Config.Outbounds[0].Settings.Vnext[0].Address = config.Host
	config.Config.Outbounds[0].Settings.Vnext[0].Users[0].Id = config.Id
	config.Config.Outbounds[0].Settings.Vnext[0].Users[0].AlterId = 0
	config.Config.Outbounds[0].StreamSettings.Network = req.Network

	config.Config.Outbounds[0].StreamSettings.WSConfig = &conf.WebSocketConfig{}
	if config.Protocol == "ws" {
		config.Config.Outbounds[0].StreamSettings.WSConfig.Path = req.Settings.WsSettings.Path
		config.Config.Outbounds[0].StreamSettings.KCPConfig = nil
		config.Config.Outbounds[0].StreamSettings.TCPConfig = nil
	}

	if config.Protocol == "kcp" {
		config.Config.Outbounds[0].StreamSettings.KCPConfig = &conf.KCPConfig{}
		config.Config.Outbounds[0].StreamSettings.WSConfig = nil
		config.Config.Outbounds[0].StreamSettings.TCPConfig = nil
		cap := uint32(req.Settings.KcpSettings.UplinkCapacity)
		config.Config.Outbounds[0].StreamSettings.KCPConfig.UpCap = cap
		cap2 := uint32(req.Settings.KcpSettings.DownlinkCapacity)
		config.Config.Outbounds[0].StreamSettings.KCPConfig.DownCap = cap2
		config.Config.Outbounds[0].StreamSettings.KCPConfig.Congestion = req.Settings.KcpSettings.Congestion
		headerRaw, _ := json.Marshal(req.Settings.KcpSettings.Header)
		var mp = make(map[string]string)
		json.Unmarshal(headerRaw, &mp)
		config.Config.Outbounds[0].StreamSettings.KCPConfig.HeaderConfig = mp
	}
	config.Vmess = config.Config.GetVmess(config.Alias)
	if err := conf.UpdateOneConfig(config); err != nil {
		sendError(ctx, errors.Wrap(err, "更新配置失败"))
		return
	}

	sendSuccess(ctx, true)
}

func del(ctx *gin.Context) {
	var req vo.UUIDConfigRequest
	if err := ctx.BindJSON(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}
	config, err := conf.GetConfigByUUID(req.UUID)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取配置失败"))
		return
	}
	if config.Status == conf.StatusStart {
		sendError(ctx, errors.New("运行中配置无法删除"))
		return
	}
	err = conf.DeleteOneConfig(req.UUID)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取配置失败"))
		return
	}
	sendSuccess(ctx, true)
}

func reload(ctx *gin.Context) {
	var req vo.UUIDConfigRequest
	if err := ctx.BindJSON(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	path, err := conf.ActiveRuntimeConfigFile(req.UUID)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "写入配置文件失败"))
		return
	}

	if err := runtime.Start(path); err != nil {
		sendError(ctx, errors.Wrap(err, "启动服务失败"))
		return
	}

	configs, err := conf.GetConfigList()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取配置文件列表失败"))
		return
	}
	lo.ForEach(configs, func(item *conf.ClientConfig, index int) {
		if item.UUID == req.UUID {
			item.Status = conf.StatusStart
		} else {
			item.Status = conf.StatusDown
		}
	})
	if err := conf.SaveConfigList(configs); err != nil {
		sendError(ctx, errors.Wrap(err, "更新配置失败"))
		return
	}
	rs, err := conf.GetRunningStatus()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "读取配置失败"))
		return
	}
	rs.RunningUUID = req.UUID
	conf.SetRunningStatus(rs)

	sendSuccess(ctx, true)
}

func getLocalConfig(ctx *gin.Context) {
	conf, err := conf.GetLocalConfig()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取本地配置失败"))
		return
	}
	sendSuccess(ctx, conf)
}

func updateLocalConfig(ctx *gin.Context) {
	var localAddress = []string{"0.0.0.0", "127.0.0.1"}

	var req conf.LocalConfig
	if err := ctx.BindJSON(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	if req.HttpPort < 1000 || req.HttpPort > 65535 {
		sendError(ctx, errors.New("http端口范围错误: 1000~65535"))
		return
	}
	if req.SocksPort < 1000 || req.SocksPort > 65535 {
		sendError(ctx, errors.New("http端口范围错误: 1000~65535"))
		return
	}

	if !lo.Contains(localAddress, req.HttpAddress) {
		sendError(ctx, errors.New("http地址错误: 127.0.0.1 , 0.0.0.0 "))
		return
	}

	if !lo.Contains(localAddress, req.SocksAddress) {
		sendError(ctx, errors.New("socks地址错误: 127.0.0.1 , 0.0.0.0 "))
		return
	}
	if err := conf.SetLocalConfig(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "写入配置失败"))
		return
	}
	sendSuccess(ctx, true)
}

func autoTest(ctx *gin.Context) {
	state := vo.ClientState{}

	lconf, err := conf.GetLocalConfig()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	if runtime.IsRunning() {
		state.IsRunning = true
	} else {
		sendSuccess(ctx, state)
		return
	}

	if lconf.SocksPort > 0 && helpers.CheckPort(lconf.SocksPort) {
		state.Socks = true
	}

	if lconf.HttpPort > 0 && helpers.CheckPort(lconf.HttpPort) {
		state.Http = true
	}

	activeConf, err := conf.GetActiveConfig()
	if err != nil {
		sendError(ctx, err)
		return
	}
	serverAddress := fmt.Sprintf("%s:%s", activeConf.Host, activeConf.Port)

	network := "tcp"
	if activeConf.Protocol == "kcp" || activeConf.Protocol == "mkcp" {
		network = "udp"
	}

	if len(serverAddress) > 0 && helpers.CheckPortbyProtocol(network, serverAddress) {
		state.ConnectToServer = true
	}

	if helpers.CheckPorxy(lconf.HttpProxy(), "https://www.google.com") {
		state.PorxyOK = true
	}

	sendSuccess(ctx, state)
}

func stop(ctx *gin.Context) {
	configs, err := conf.GetConfigList()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "获取配置文件列表失败"))
		return
	}
	lo.ForEach(configs, func(item *conf.ClientConfig, index int) {
		if item.Status == conf.StatusStart {
			item.Status = conf.StatusDown
		}
	})
	runtime.Stop()
	if err := conf.SaveConfigList(configs); err != nil {
		sendError(ctx, errors.Wrap(err, "写入配置失败"))
		return
	}
	sendSuccess(ctx, true)
}

func qrCode(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		return
	}
	data, _ := qrcode.Encode(code, qrcode.High, 300)
	ctx.Data(200, "image/png", data)
	return
}

type VmessImportRequest struct {
	Vmess string `json:"vmess"`
}

func importVmess(ctx *gin.Context) {
	var req VmessImportRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	if req.Vmess == "" {
		sendError(ctx, errors.New("参数不能为空"))
		return
	}

	link, err := vmess.FromVmessLink(req.Vmess)
	if err != nil {
		sendError(ctx, err)
		return
	}
	var clientConfig conf.ClientConfigV2ray
	var defaultConfigTemplate = conf.ClientDefaultConfig
	err = json.Unmarshal([]byte(defaultConfigTemplate), &clientConfig)
	if err != nil {
		sendError(ctx, err)
		return
	}

	clientConfig.Outbounds[0].Settings.Vnext[0].Users[0].Id = link.Id
	clientConfig.Outbounds[0].Settings.Vnext[0].Users[0].Id = link.Id
	clientConfig.Outbounds[0].Settings.Vnext[0].Address = link.Host
	clientConfig.Outbounds[0].Settings.Vnext[0].Port = str2Int(link.Port)
	clientConfig.Outbounds[0].StreamSettings.Network = link.Network
	clientConfig.Outbounds[0].StreamSettings.Security = link.Tls

	switch link.Network {
	case "ws":
		clientConfig.Outbounds[0].StreamSettings.WSConfig = &conf.WebSocketConfig{}
		if link.Path != "" && link.Host != "" {
			clientConfig.Outbounds[0].StreamSettings.WSConfig = &conf.WebSocketConfig{
				Path:  link.Path,
				Path2: "",
				Headers: map[string]string{
					"Host": link.Host,
				},
			}
		}
	case "kcp":
		clientConfig.Outbounds[0].StreamSettings.KCPConfig = &conf.KCPConfig{}
		var up uint32 = 5
		var down uint32 = 100
		header := map[string]string{
			"type": "none",
		}
		if link.Type != "" {
			header["type"] = link.Type
		}
		clientConfig.Outbounds[0].StreamSettings.KCPConfig = &conf.KCPConfig{
			UpCap:        up,
			DownCap:      down,
			HeaderConfig: header,
		}
	}

	clientConfig.UUID = uuid.New().String()

	cc := &conf.ClientConfig{
		UUID:     clientConfig.UUID,
		Config:   &clientConfig,
		Alias:    link.Name,
		Vmess:    link.ToVmessLink(),
		Host:     link.Address,
		Port:     link.Port,
		Protocol: link.Network,
		Id:       link.Id,
		Ts:       time.Now().Unix(),
	}

	if err := conf.CreateOneConfig(cc); err != nil {
		sendError(ctx, errors.Wrap(err, "保存配置失败"))
		return
	}
	sendSuccess(ctx, true)
}

type ConfigJsonRequest struct {
	Uuid string `json:"uuid"`
}

func configJson(ctx *gin.Context) {
	var req ConfigJsonRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	config, err := conf.GetConfigByUUID(req.Uuid)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "查找配置失败"))
		return
	}

	localConfig, err := conf.GetLocalConfig()
	if err != nil {
		sendError(ctx, errors.Wrap(err, "查找配置失败"))
		return
	}
	inBounds := config.Config.Inbounds
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
	config.Config.Inbounds = inBounds
	config.Config.Log.Access = filepath.Join(conf.GetDefaultConfigDirectory(), "access.log")
	config.Config.Log.Error = filepath.Join(conf.GetDefaultConfigDirectory(), "error.log")

	data, err := json.MarshalIndent(config.Config, "", "\t")
	if err != nil {
		sendError(ctx, errors.Wrap(err, "序列号配置失败"))
		return
	}

	sendSuccess(ctx, string(data))
}

func str2Int(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

type ConfigImportRequest struct {
	Config string `json:"config"`
}

func configImport(ctx *gin.Context) {
	var req ConfigImportRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}

	clientV2ray, err := conf.ClientConfigV2rayFromString(req.Config)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "配置解析失败"))
		return
	}

	id := uuid.New().String()
	clientV2ray.UUID = id
	name := "导入配置-" + time.Now().Format(`2006-01-02 15:04:05`)
	vmessLinkString := clientV2ray.GetVmess(name)
	vmessLink, err := vmess.FromVmessLink(vmessLinkString)
	if err != nil {
		sendError(ctx, errors.Wrap(err, "配置vmess失败"))
		return
	}
	clientConfig := &conf.ClientConfig{
		UUID:     uuid.New().String(),
		Config:   clientV2ray,
		Alias:    name,
		Vmess:    clientV2ray.GetVmess(name),
		Host:     vmessLink.Host,
		Port:     vmessLink.Port,
		Protocol: vmessLink.Network,
		Id:       vmessLink.Id,
		Status:   0,
		Ts:       time.Now().Unix(),
	}

	if err := conf.CreateOneConfig(clientConfig); err != nil {
		sendError(ctx, errors.Wrap(err, "保存配置失败"))
		return
	}

	sendSuccess(ctx, true)
}

func shutdown(ctx *gin.Context) {
	sendSuccess(ctx, true)
	menu.Shutdown()
}

func systemProxyStatus(ctx *gin.Context) {
	status := system.GetProxyStatus()
	mode := system.GetMode()
	sendSuccess(ctx, map[string]interface{}{
		"status":     status,
		"mode":       mode,
		"pacAddress": GetPacAddress(),
	})
}

func Pacjs(ctx *gin.Context) {
	ctx.Header("Content-type", "application/x-ns-proxy-autoconfig")
	ctx.String(200, conf.ParsePacJS())
}

type SetProxyRequest struct {
	Mode int `json:"mode"`
}

func setProxy(ctx *gin.Context) {
	var req SetProxyRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, errors.Wrap(err, "参数绑定失败"))
		return
	}
	var err error
	switch req.Mode {
	case system.ModeNone:
		err = system.SetNone()
	case system.ModePac:
		err = system.SetPac(GetPacAddress())
	case system.ModeGlobal:
		err = system.SetGlobal()
	}
	if err != nil {
		sendError(ctx, errors.Wrap(err, "操作失败"))
		return
	}
	sendSuccess(ctx, true)
}
