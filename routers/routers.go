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

		//userInfo
		r.POST("/"+conEntrance+"/*action", admin.UserInfo)
	}

	return r
}
