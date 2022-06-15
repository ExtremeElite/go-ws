package util

import "encoding/json"

const (
	ContentType          = "Content-Type"
	AppJson              = "application/json"
	HelloWorld           = `<h1>欢迎来到Gorouting即时通讯服务</h1>`
	TimeOut              = "请求超时"
	ReadConnectClosed    = "读连接关闭"
	WriteConnectClosed   = "写连接关闭"
	ConnectClosed        = "连接已经关闭"
	MysqlTcpConnect      = `%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`
	MessageHeaderSuccess = `[success]`
	MessageHeaderFailed  = `[failed]`
	NotFound             = `{"code":404,"msg":"数据错误","data":""}`
	MethodNotAllowed     = `{"code":405,"msg":"请求方式错误","data":""}`
	HostNotAllowed       = `{"code":406,"msg":"请求域名错误","data":""}`
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (response Response) Json(msg string, code int, data interface{}) []byte {
	response.Msg = msg
	response.Code = code
	response.Data = data
	res, err := json.Marshal(response)
	if err != nil {
		return []byte(NotFound)
	}
	return res
}
