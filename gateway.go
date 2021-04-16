package main

import (
	"go-gateway/inc"
	"go-gateway/routers"
)

func main() {
	//初始化路由
	r := routers.InitRouter()

	// 监听端口，默认在8080
	port := inc.Cfg.MustValue("http", "port", "8080")
	r.Run(":" + port)

}
