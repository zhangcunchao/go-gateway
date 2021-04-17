package main

import (
	"go-gateway/build"
	"go-gateway/debug"
	"go-gateway/inc"
	"go-gateway/routers"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			debug.DebugPrint("error：%s", r)
		}
	}()

	var args = os.Args

	if 2 > len(args) {
		args = append(args, "")
	}

	switch a := args[1]; {
	case a == "init":
		if 3 > len(args) {
			args = append(args, "")
		}
		build.InitDb(args[2])
	default:
		inc.InitGatway()
		r := routers.InitRouter()
		// 监听端口，默认在8080
		port := inc.Cfg.MustValue("http", "port", "8080")
		r.Run(":" + port)
	}

}
