package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v2ray_panel_plus "github.com/make-money-fast/v2ray-panel-plus"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/menu"
	"github.com/phayes/freeport"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ListenAddress = ""

func StartHttpServer() {
	g := gin.Default()
	port := getListenPort()

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

	g.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/static/client")
	})
	g.GET("/proxy.pac", Pacjs)
	g.GET("/qrcode", qrCode)

	ListenAddress = "http://127.0.0.1:" + port
	// html
	fmt.Println("Servering at on ", ListenAddress)
	menu.UIAddress = ListenAddress
	g.Run(fmt.Sprintf(":%s", port))
}

func getListenPort() string {
	port := os.Getenv("V2RAY_PANEL_LISTEN_ADDRESS")
	if port == "" {
		p, _ := freeport.GetFreePort()
		port = fmt.Sprintf("%d", p)
	}
	return port
}

func GetPacAddress() string {
	if ListenAddress == "" {
		return ""
	}
	return ListenAddress + "/proxy.pac?t=" + strconv.Itoa(int(time.Now().Unix()))
}
