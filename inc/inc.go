package inc

import (
	"os"
	"runtime"

	"go-gateway/debug"

	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
)

var (
	Cfg *goconfig.ConfigFile
)

//主动初始化
func Init() {

	defer func() {
		if r := recover(); r != nil {
			debug.DebugPrint("捕获到错误：%s\n", r)
			os.Exit(2)
		}
	}()

	//初始化conf.ini
	var err error
	Cfg, err = goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic("conf.ini不存在")
	}

	//多核配置
	numCPU := Cfg.MustInt("http", "numCPU", 0)
	if numCPU == 0 {
		numCPU = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(numCPU)
	debug.DebugPrint("run cpuNum %d", numCPU)

	//gin运行模式
	DebugModel := Cfg.MustValue("http", "DebugModel", "debug")
	gin.SetMode(DebugModel)
}
