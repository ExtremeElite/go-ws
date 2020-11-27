/**
 * @date:2020\11\27 0027 22:02
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"ws/core"
)

const (
	Conversation  =1
	Login         =2
	Logout        =3
	GetOnlineInfo =4

)
type wsPipeLineFn func([]byte,*core.Connection) error
type Response struct {
	Code uint8 `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
//推送格式
type PushData struct {
	EventType int32
	Device []string
	Data string
}

func (response Response)Json() []byte{
	data,err:=json.Marshal(response)
	if err!=nil {
		return []byte(`["code":404,"msg":"数据错误"]`)
	}
	return data
}