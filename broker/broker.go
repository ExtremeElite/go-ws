/**
 * @date:2020\11\27 0027 22:02
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
//推送格式
type PushData struct {
	EventType int32 `json:"event_type"`
	PublishAccount []string `json:"publish_account"`
	Data interface{} `json:"data"`
}
func (response Response)Json(msg string,code int,data interface{}) []byte{
	response.Msg=msg
	response.Code=code
	response.Data=data
	res,err:=json.Marshal(response)
	if err!=nil {
		return []byte(`["code":404,"msg":"数据错误","data":""]`)
	}
	return res
}
//格式转换
func (pushData PushData) ConversionJson() string{
	data:=pushData.Data
	dataType:=strings.ToLower(fmt.Sprintf("%s",reflect.TypeOf(data).Kind()))
	//字符串
	if dataType=="string" {
		return data.(string)
	}
	//数字
	if dataType=="float64" {
		return strings.TrimRight(strconv.FormatFloat(data.(float64), 'E', -1, 64),`E+00`)
	}
	//对象或者数组
	if dataType=="map" || dataType=="slice" {
		result,err:=json.Marshal(data)
		if err!=nil {
			println(err.Error())
			return ""
		}
		return string(result)
	}
	return ""
}