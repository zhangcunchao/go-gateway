package admin

import (
	"encoding/base64"
	"encoding/json"
	"go-gateway/debug"
	"go-gateway/exception"
	"go-gateway/inc"
	"go-gateway/model"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const COOKIE_TIMEOUT = 86400 * 1

type adminAuth struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	LoginDate time.Time `json:"login_date"`
}

var UserLoginInfo adminAuth

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断登录
		if cookie, err := c.Request.Cookie("adminauth"); err == nil {
			value := cookie.Value
			value, _ = url.PathUnescape(value)
			//value = strings.Replace(value, "%3D", "=", -1)
			e, err := base64.URLEncoding.DecodeString(value)
			//debug.DebugPrint("sub", err)
			if err == nil {
				//解密
				res := inc.DesDecrypt_CBC(e)
				err := json.Unmarshal(res, &UserLoginInfo)
				if err == nil && UserLoginInfo.Id > 0 {
					//cookie 有效期校验
					now := time.Now()
					time := now.Sub(UserLoginInfo.LoginDate)
					debug.DebugPrint("sub", time)
					if time > COOKIE_TIMEOUT {
						//账户状态等校验
						admin, row := model.GetFirstAdmin("id = ?", UserLoginInfo.Id)
						if row > 0 {
							if admin.Status == 1 {
								c.Next()
								return
							} else {
								Return(exception.COOD_FAIL_10003, "账户被禁用", "", c)
								c.Abort()
								return
							}

						}
					}

				}
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

//登录账户信息
func UserInfo(c *gin.Context) {

	Return(exception.COOD_SUCCESS, "调用成功", UserLoginInfo, c)
}

//管理员列表
func UserList(c *gin.Context) {

	adminList, _ := model.GetListAdmin("1 = ?", 1)

	Return(exception.COOD_SUCCESS, "调用成功", adminList, c)
}

func Login(c *gin.Context) {
	var loginInfo model.AdminLoginInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		Return(exception.COOD_FAIL, "请输入用户名、密码", nil, c)
		//Return(exception.COOD_FAIL, err.Error(), nil, c)
		return
	}

	var admin model.Admin
	var err error
	if admin, err = model.DoLogin(loginInfo); err != nil {
		Return(exception.COOD_FAIL, err.Error(), nil, c)
		return
	}

	auth := adminAuth{Id: admin.ID, Name: admin.Name, LoginDate: time.Now()}
	b, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}
	//debug.DebugPrint("adminauth1", string(b))
	result := inc.DesEncrypt_CBC(b)

	f := base64.URLEncoding.EncodeToString([]byte(result))
	//f := base64.StdEncoding.EncodeToString(result)

	conEntrance := inc.Cfg.MustValue("http", "ConEntrance", "")
	//按path设置cookie
	debug.DebugPrint("cookie", f)
	c.SetCookie("adminauth", f, COOKIE_TIMEOUT, "/"+conEntrance, "", false, false)
	Return(exception.COOD_SUCCESS, "调用成功", nil, c)

}

func Return(code int, msg interface{}, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, exception.APIException{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
