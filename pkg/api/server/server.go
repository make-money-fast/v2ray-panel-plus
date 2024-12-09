package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/api/vo"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/ufw"
	"github.com/samber/lo"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func sendError(ctx *gin.Context, err error) {
	ctx.JSON(200, &Response{
		Code: -1,
		Msg:  err.Error(),
	})
}

func sendSuccess(ctx *gin.Context, data any) {
	ctx.JSON(200, &Response{
		Code: 0,
		Data: data,
	})
}

func Reload(ctx *gin.Context) {
	path, err := conf.ActiveServerRuntimeConfig()
	if err != nil {
		sendError(ctx, err)
		return
	}
	if err := runtime.Reload(path); err != nil {
		sendError(ctx, err)
		return
	}
	ufw.ActiveUFW()
	sendSuccess(ctx, nil)
	return
}

func RuntimeConfig(ctx *gin.Context) {
	conf, err := conf.MergeServerConfig()
	if err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, conf)
	return
}

func ListConfig(ctx *gin.Context) {
	data, err := conf.GetServerConfigList()
	if err != nil {
		return
	}
	sendSuccess(ctx, data)
}

func EditConfig(ctx *gin.Context) {
	var req vo.EditConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	list, err := conf.GetServerConfigList()
	if err != nil {
		sendError(ctx, err)
		return
	}
	lo.ForEach(list, func(item *conf.ServerConfig, index int) {
		if item.UUID == req.UUID {
			list[index] = req.ToConfig()
		}
	})
	if err := conf.SaveServerConfigList(list); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, nil)
}

func AddConfig(ctx *gin.Context) {
	var req vo.EditConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	list, err := conf.GetServerConfigList()
	if err != nil {
		sendError(ctx, err)
		return
	}
	item := req.ToConfig()
	list = append(list, item)
	if err := conf.SaveServerConfigList(list); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, nil)
}

func DelConfig(ctx *gin.Context) {
	var req vo.DelConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	if err := conf.DeleteOneServerConfig(req.UUID); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, nil)
}
