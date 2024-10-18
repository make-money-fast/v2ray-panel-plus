package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
	"net/http"
	"os"
	"strconv"
	"time"
	v2ray_panel_plus "v2ray-panel-plus"
	"v2ray-panel-plus/pkg/menu"
)

var ListenAddress = ""

func StartHttpServer() {
	g := gin.Default()
	listenAddress := getListenAddress()

	if os.Getenv("DEBUG") == "1" {
		g.Static("/static", "static")
	} else {
		gin.SetMode(gin.ReleaseMode)
		g.Any("/static/*action", func(ctx *gin.Context) {
			http.FileServer(http.FS(v2ray_panel_plus.StaticFS)).ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	api := g.Group("/api")
	{
		api.POST("/init", initDefaultConfig)
		api.POST("/list", getConfigList)
		api.POST("/edit", editConfig)
		api.POST("/del", del)
		api.POST("/reload", reload)

		api.POST("/getLocalConfig", getLocalConfig)
		api.POST("/updateLocalConfig", updateLocalConfig)

		api.POST("/autoTest", autoTest)
		api.POST("/stop", stop)
		api.POST("/importVmess", importVmess)
		api.POST("/configJSON", configJson)
		api.POST("/config-import", configImport)
		api.POST("/shutdown", shutdown)
		api.POST("/systemProxyStatus", systemProxyStatus)
		api.POST("/setProxy", setProxy)
	}

	g.GET("/proxy.pac", Pacjs)
	g.GET("/qrcode", qrCode)

	ListenAddress = "http://" + listenAddress
	// html
	fmt.Println("Listening on ", ListenAddress)
	uiAddress := "http://" + listenAddress + "/static/client"
	menu.UIAddress = uiAddress

	g.Run(listenAddress)
}

func getListenAddress() string {
	port := os.Getenv("V2RAY_PANEL_LISTEN_ADDRESS")
	if port == "" {
		p, _ := freeport.GetFreePort()
		port = fmt.Sprintf("%d", p)
	}
	return fmt.Sprintf("0.0.0.0:%s", port)
}

func GetPacAddress() string {
	if ListenAddress == "" {
		return ""
	}
	return ListenAddress + "/proxy.pac?t=" + strconv.Itoa(int(time.Now().Unix()))
}
