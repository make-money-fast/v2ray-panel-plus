package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/api/server"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/conf"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/runtime/client"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var daemon bool

func main() {
	if len(os.Args) == 1 {
		daemon = true
		fmt.Println("server start daemon. listen http interface: http://127.0.0.1:7810")
		startDaemon()
		return
	}
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	rootCmd = cli.Command{}
	ch      chan struct{}
)

func init() {
	rootCmd.Commands = []*cli.Command{
		{
			Name:      "start",
			Usage:     "启动v2ray服务",
			UsageText: "",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionStart(ctx, command)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "uuid",
					Value: "",
				},
			},
		},
		{
			Name:  "list",
			Usage: "列出服务列表",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionList(ctx, command)
			},
			Flags: []cli.Flag{},
		},
		{
			Name:  "del",
			Usage: "列出服务列表",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionDel(ctx, command)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "uuid",
					Value: "",
				},
			},
		},
		{
			Name:      "add",
			Usage:     "增加服务端配置",
			UsageText: "",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionAdd(ctx, command)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file",
					Value: "",
					Usage: "增加的文件名称",
				},
				&cli.StringFlag{
					Name:  "alias",
					Value: "",
					Usage: "别名",
				},
			},
		},
		{
			Name:  "stop",
			Usage: "停止v2ray服务",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionStop(ctx, command)
			},
		},
		{
			Name:  "reload",
			Usage: "重启v2ray服务",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionStop(ctx, command)
			},
		},
		{
			Name:  "link",
			Usage: "vmess地址",
			Action: func(ctx context.Context, command *cli.Command) error {
				return actionLink(ctx, command)
			},
		},
	}
}

func init() {
	conf.AsServer()
	conf.InitDefaultConfigFile()
	conf.InitTemplateFile()
}

func actionList(ctx context.Context, command *cli.Command) error {
	serverConfigList, err := conf.GetServerConfigList()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetTitle("Server List")

	t.AppendHeader(table.Row{"UUID", "Alias", "Network", "Port"})
	lo.ForEach(serverConfigList, func(item *conf.ServerConfig, index int) {
		t.AppendRow(table.Row{item.UUID, item.Alias, item.Protocol, item.Port})
	})
	fmt.Println(t.Render())
	return nil
}

func actionAdd(ctx context.Context, command *cli.Command) error {
	file := command.String("file")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var serverConfig conf.ServerConfig
	err = json.Unmarshal(data, &serverConfig)
	if err != nil {
		return err
	}
	alias := command.String("alias")
	if alias != "" {
		serverConfig.Alias = alias
	}
	if serverConfig.Port == "" {
		serverConfig.Port = strconv.Itoa(serverConfig.Config.Port)
	}
	if serverConfig.UUID == "" {
		serverConfig.UUID = uuid.New().String()
	}
	if serverConfig.Protocol == "" {
		serverConfig.Protocol = serverConfig.Config.StreamSettings.Network
	}
	if serverConfig.Id == "" {
		if len(serverConfig.Config.Settings.Clients) == 0 {
			return errors.New("client id is empty")
		}
		serverConfig.Id = serverConfig.Config.Settings.Clients[0].Id
	}
	if serverConfig.Alias == "" {
		serverConfig.Alias = fmt.Sprintf("%s:%s", serverConfig.Protocol, serverConfig.Port)
	}
	if serverConfig.Type == "" && serverConfig.Config.StreamSettings != nil && serverConfig.Config.StreamSettings.KCPConfig != nil && serverConfig.Config.StreamSettings.KCPConfig.HeaderConfig != nil {
		serverConfig.Type = serverConfig.Config.StreamSettings.KCPConfig.HeaderConfig["type"]
	}
	if err := conf.CreateOneServerConfig(&serverConfig); err != nil {
		return err
	}
	fmt.Println("server add new configure: uuid=", serverConfig.UUID)
	fmt.Println("protocol: ", serverConfig.Protocol)
	fmt.Println("port: ", serverConfig.Port)
	fmt.Println("id: ", serverConfig.Id)
	fmt.Println("vmess Link: ", serverConfig.BuildVmess())
	return nil
}

func actionStart(ctx context.Context, command *cli.Command) error {
	if daemon {
		path, err := conf.ActiveServerRuntimeConfig()
		if err != nil {
			return err
		}
		err = client.Start(path)
		if err != nil {
			return err
		}
		fmt.Println("server start ok")
		return nil
	}
	rawRequest(os.Args[1:])
	return nil
}

func actionStop(ctx context.Context, command *cli.Command) error {
	if daemon {
		client.Stop()
		fmt.Println("server stop ok")
		return nil
	}
	rawRequest(os.Args[1:])
	return nil
}

func actionReload(ctx context.Context, command *cli.Command) error {
	if daemon {
		path, err := conf.ActiveServerRuntimeConfig()
		if err != nil {
			return err
		}
		if err := client.Reload(path); err != nil {
			return err
		}
		fmt.Println("server reload ok")
		return nil
	}
	rawRequest(os.Args[1:])
	return nil
}

func actionDel(ctx context.Context, command *cli.Command) error {
	id := command.String("uuid")
	if err := conf.DeleteOneServerConfig(id); err != nil {
		return err
	}
	fmt.Println("deleted: ", id)
	return nil
}

type JsonRequest struct {
	Args []string `json:"args"`
}

func startDaemon() {
	go func() {
		server.RunServer()
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		var jsonReq JsonRequest
		err = json.Unmarshal(data, &jsonReq)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = rootCmd.Run(r.Context(), append([]string{os.Args[0]}, jsonReq.Args...))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("ok"))
	})
	http.ListenAndServe("127.0.0.1:7810", nil)
}

func rawRequest(args []string) {
	fmt.Println("args: ", args)
	var req JsonRequest
	req.Args = args
	data, _ := json.Marshal(req)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:7810", bytes.NewReader(data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Add("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
	if string(data) == "ok" {
		fmt.Println("handle success")
	}
	return
}

func actionLink(ctx context.Context, command *cli.Command) error {
	list, err := conf.GetServerConfigList()
	if err != nil {
		return err
	}
	lo.ForEach(list, func(item *conf.ServerConfig, index int) {
		uuid := command.String("uuid")
		if uuid != "" && item.UUID != uuid {
			return
		}
		fmt.Printf("%s: %s\n\n", item.Alias, item.BuildVmess())
	})

	return nil
}
