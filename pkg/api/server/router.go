package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	g := gin.Default()

	api := g.Group("/api")
	{
		api.POST("/reload", Reload)
		api.POST("/runtime", RuntimeConfig)
		api.POST("/list", ListConfig)
		api.POST("/edit", EditConfig)
		api.POST("/add", AddConfig)
		api.POST("/del", DelConfig)
	}

	fmt.Println("api server running at: http://localhost:7677")
	g.Run(":7677")
}
