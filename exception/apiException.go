package exception

const COOD_SUCCESS = 10000   //成功
const COOD_FAIL = 10001      //失败
const COOD_NOT_LOGIN = 10002 //未登陆

// api错误的结构体
type APIException struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
