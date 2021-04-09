package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOD_SUCCESS = 10000
const COOD_FAIL = 10001

func UserInfo(c *gin.Context) {

	data := map[string]string{"userName": "fuck", "id": "123"}

	Return(COOD_SUCCESS, "调用成功", data, c)
}

func Return(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
