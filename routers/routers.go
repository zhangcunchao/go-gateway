package routers

import (
	"net/http"
	"os"
	"path"
	"time"

	"go-gateway/controller/admin"
	"go-gateway/debug"
	"go-gateway/exception"
	"go-gateway/inc"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
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

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(logerMiddleware(), Recover)

	gwEntrance := inc.Cfg.MustValue("http", "GwEntrance", "api")
	r.Any("/"+gwEntrance+"/*action", entrance)

	debug.DebugPrint("gateway entrance %s", gwEntrance)
	//管理后台入口
	conEntrance := inc.Cfg.MustValue("http", "ConEntrance", "")
	if conEntrance != "" {
		r.StaticFS("/"+conEntrance, http.Dir("./page"))

		r.POST("/"+conEntrance+"/login", admin.Login)

		av1 := r.Group(conEntrance).Use(admin.AuthMiddleWare())
		//登录用户信息
		av1.POST("userInfo", admin.UserInfo)
		//管理账户列表
		av1.POST("userList", admin.UserList)

	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": exception.COOD_FAIL_10004,
			"msg":  "404 NotFound",
			"data": nil,
		})
	})

	return r
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {

			c.JSON(http.StatusOK, gin.H{
				"code": exception.COOD_FAIL_10005,
				"msg":  r,
				"data": nil,
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

func logerMiddleware() gin.HandlerFunc {

	// 日志文件
	logPath := inc.Cfg.MustValue("http", "logPath", "storage/logs/access") // 日志打印到指定的目录
	logMaxAge := inc.Cfg.MustInt("http", "logMaxAge", 15)

	fileName := path.Join(logPath, "access.log")

	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		debug.DebugPrint("gg-gin日志异常", err)
	}
	// 实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//logger.SetFormatter(&logrus.JSONFormatter{})
	//设置输出
	logger.Out = src

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y-%m-%d",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		debug.DebugPrint("gg-gin日志异常", err)
	}

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()

	}
}
