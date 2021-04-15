package routers

import (
	"net/http"

	"go-gateway/controller/admin"
	"go-gateway/debug"
	"go-gateway/inc"

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
)

//初始化
func init() {
	//初始化路由集合 todo 判断设置了api才启用
	proxyMap["/abc"] = Proxy{Remark: "Remark", Prefix: "ab"}
}

//网关入口
func entrance(c *gin.Context) {
	name := c.Param("action")

	if val, ok := proxyMap[name]; ok {
		//路由是否存在
		//fmt.Println(val)
		debug.DebugPrint("后端接口%s存在%s", name, val)
	} else {
		// fmt.Println(ok)
		// fmt.Printf("%s不存在\n", name)
		debug.DebugPrint("后端接口%s不存在\n", name)
	}
	c.String(http.StatusOK, "hello World! %s", name)
}

func abc(c *gin.Context) {
	name := c.FullPath()
	c.String(http.StatusOK, "hello World! %s", name)
}

func InitRouter() *gin.Engine {
	r := gin.New()

	gwEntrance := inc.Cfg.MustValue("http", "GwEntrance", "api")
	r.Any("/"+gwEntrance+"/*action", entrance)

	test := "/abcd/:id"

	r.Any(test, abc)

	debug.DebugPrint("gateway entrance %s", gwEntrance)
	//管理后台入口
	conEntrance := inc.Cfg.MustValue("http", "ConEntrance", "")
	if conEntrance != "" {
		r.StaticFS("/"+conEntrance, http.Dir("./page"))

		r.POST("/"+conEntrance+"/login", admin.Login)

		av1 := r.Group(conEntrance).Use(admin.AuthMiddleWare())
		av1.POST("userinfo", admin.UserInfo)

	}

	return r
}
