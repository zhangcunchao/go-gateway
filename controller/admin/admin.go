package admin

import (
	"go-gateway/exception"
	"go-gateway/inc"
	"go-gateway/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		Return(exception.COOD_NOT_LOGIN, "请先登录", "", c)
		c.Abort()
	}
}

func UserInfo(c *gin.Context) {

	admin, _ := model.GetFirstAdmin("id = ? and status = ?", 1, 1)

	Return(exception.COOD_SUCCESS, "调用成功", admin, c)
}

func Login(c *gin.Context) {
	var loginInfo model.AdminLoginInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		Return(exception.COOD_FAIL, "请输入用户名、密码", nil, c)
		//Return(exception.COOD_FAIL, err.Error(), nil, c)
		return
	}

	if _, err := model.DoLogin(loginInfo); err != nil {
		Return(exception.COOD_FAIL, err.Error(), nil, c)
		return
	}

	conEntrance := inc.Cfg.MustValue("http", "ConEntrance", "")
	//按path设置cookie
	c.SetCookie("adminauth", "onion", COOKIE_TIMEOUT, "/"+conEntrance, "", false, false)
	Return(exception.COOD_SUCCESS, "调用成功", nil, c)

}

func Return(code int, msg interface{}, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, exception.APIException{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
