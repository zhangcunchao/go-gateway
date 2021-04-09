package routers

import (
	"net/http"

	"go-gateway/controller/admin"

	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
)

func InitRouter(Cfg *goconfig.ConfigFile) *gin.Engine {
	r := gin.New()
	//管理后台入口
	conEntrance := Cfg.MustValue("http", "ConEntrance", "")
	if conEntrance != "" {
		r.StaticFS("/"+conEntrance, http.Dir("./page"))

		r.POST("/"+conEntrance+"/login", admin.Login)
		r.Use(admin.AuthMiddleWare())
		{
			//userInfo
			r.POST("/"+conEntrance+"/userinfo", admin.UserInfo)
		}
	}

	return r
}
