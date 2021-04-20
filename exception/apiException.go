package exception

const COOD_SUCCESS = 10000    //成功
const COOD_FAIL = 10001       //失败
const COOD_NOT_LOGIN = 10002  //未登陆
const COOD_FAIL_10003 = 10003 //账户禁用
const COOD_FAIL_10004 = 10004 //路由不存在
const COOD_FAIL_10005 = 10005 //运行异常

// api错误的结构体
type APIException struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
