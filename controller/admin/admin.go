package admin

import (
	"go-gateway/inc"
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOD_SUCCESS = 10000
const COOD_FAIL = 10001
const COOD_NOT_LOGIN = 10002

const COOKIE_TIMEOUT = 6000

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断登录
		if cookie, err := c.Request.Cookie("adminauth"); err == nil {
			value := cookie.Value
			if value == "onion" {
				c.Next()
				return
			}
		}
		// if url := c.Request.URL.String(); url == "/login" {
		// 	c.Next()
		// 	return
		// }
		//判断用户状态权限
		Return(COOD_NOT_LOGIN, "请先登录", "", c)
		c.Abort()
	}
}

func UserInfo(c *gin.Context) {

	data := map[string]string{"userName": "fuck", "id": "123"}

	Return(COOD_SUCCESS, "调用成功", data, c)
}

func Login(c *gin.Context) {
	loginInfo := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&loginInfo)

	conEntrance := inc.Cfg.MustValue("http", "ConEntrance", "")
	//按path设置cookie
	c.SetCookie("adminauth", "onion", COOKIE_TIMEOUT, "/"+conEntrance, "", false, false)

	Return(COOD_SUCCESS, "调用成功", "", c)
}

func Return(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
