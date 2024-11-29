package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/api/vo"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func listServerConfig(ctx *gin.Context) {
	var req vo.RequestContext
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	var resp Response[[]conf.ServerConfig]
	if err := doServerRequest(req.Url("/list"), nil, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].UUID < resp.Data[j].UUID
	})
	sendSuccess(ctx, resp.Data)
}

func runtimeServerConfig(ctx *gin.Context) {
	var req vo.RequestContext
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	var resp Response[conf.ServerV2rayConfig]
	if err := doServerRequest(req.Url("/runtime"), nil, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	data, _ := json.MarshalIndent(resp.Data, "", "\t")
	sendSuccess(ctx, string(data))
}

func editServerConfig(ctx *gin.Context) {
	var req vo.ClientEditConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	var resp Response[any]
	if err := doServerRequest(req.RequestContext.Url("/edit"), req.EditConfigRequest, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, true)
}

func deleteServerConfig(ctx *gin.Context) {
	var req vo.ClientDelConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	var resp Response[any]
	if err := doServerRequest(req.RequestContext.Url("/del"), req.DelConfigRequest, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, true)
}

func reloadServer(ctx *gin.Context) {
	var req vo.RequestContext
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	var resp Response[any]
	if err := doServerRequest(req.Url("/reload"), nil, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, true)
}

func getUUID(ctx *gin.Context) {
	id := uuid.New().String()
	sendSuccess(ctx, id)
}

func addServerConfig(ctx *gin.Context) {
	var req vo.ClientEditConfigRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	req.UUID = uuid.New().String()
	var resp Response[any]
	if err := doServerRequest(req.RequestContext.Url("/add"), req.EditConfigRequest, &resp); err != nil {
		sendError(ctx, err)
		return
	}
	sendSuccess(ctx, true)
}

func doServerRequest[T any](url string, body interface{}, response *Response[T]) error {
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	} else {
		reader = strings.NewReader("{}")
	}
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", vo.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}
	return nil
}
