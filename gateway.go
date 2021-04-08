package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
)

type Proxy struct {
	Remark        string //描述
	Prefix        string //转发的前缀判断
	Upstream      string //后端 nginx 地址或者ip地址
	RewritePrefix string //重写
}

var (
	proxyMap = make(map[string]Proxy)
	Cfg      *goconfig.ConfigFile
)

//初始化
func init() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到错误：%s\n", r)
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
	debugPrint("run cpuNum %d", numCPU)

	//gin运行模式
	DebugModel := Cfg.MustValue("http", "DebugModel", "debug")
	gin.SetMode(DebugModel)

	//初始化路由集合 todo
	proxyMap["/abc"] = Proxy{Remark: "Remark", Prefix: "ab"}

}

//网关入口
func entrance(c *gin.Context) {
	name := c.Param("action")

	if val, ok := proxyMap[name]; ok {
		//路由是否存在
		//fmt.Println(val)
		debugPrint("后端接口%s存在%s", name)
	} else {
		// fmt.Println(ok)
		// fmt.Printf("%s不存在\n", name)
		debugPrint("后端接口%s不存在\n", name)
	}
	c.String(http.StatusOK, "hello World! %s", name)
}

func main() {

	//todo 分布式，控制台分离,统一返回（可配置，返回json,xml,html）

	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	// r.Any("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "hello World!"+c.ClientIP())
	// })

	gwEntrance := Cfg.MustValue("http", "GwEntrance", "api")
	r.Any("/"+gwEntrance+"/*action", entrance)
	debugPrint("gateway entrance %s", gwEntrance)

	// 1.json
	// r.GET("/someJSON", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "someJSON", "status": 200})
	// })
	// 2. 结构体响应
	// r.GET("/someStruct", func(c *gin.Context) {
	// 	type S struct {
	// 		Type string `json:"type"`
	// 	}
	// 	var msg struct {
	// 		Name    string `json:"name"`
	// 		Message string `json:"message"`
	// 		Tty     S      `json:"tty"`
	// 	}
	// 	//type I map[string]map[string]interface{}

	// 	msg.Name = "root"
	// 	msg.Message = "message"
	// 	msg.Tty = S{Type: "abcd"}

	// 	//var m = I{"abc": {"a": "中文"}, "f": {"d": 3}}
	// 	//k := "abc"

	// 	// if val, ok := m[k]; ok {
	// 	// 	//路由是否存在
	// 	// 	fmt.Println(val)
	// 	// 	c.String(http.StatusOK, "%s=%s", k, val)
	// 	// } else {
	// 	// 	fmt.Println(ok)
	// 	// 	c.String(http.StatusOK, "%s不存在", k)
	// 	// }

	// 	c.JSON(200, msg)
	// })
	// 3.XML
	// r.GET("/someXML", func(c *gin.Context) {
	// 	c.XML(200, gin.H{"message": "abc"})
	// })
	// 4.YAML响应
	// r.GET("/someYAML", func(c *gin.Context) {
	// 	c.YAML(200, gin.H{"name": "zhangsan"})
	// })

	// 监听端口，默认在8080
	port := Cfg.MustValue("http", "port", "8080")
	r.Run(":" + port)
}
