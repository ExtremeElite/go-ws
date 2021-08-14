/**
 * @date:2020\11\27 0027 22:02
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"go/types"
	"log"
	"strconv"
	"strings"
	"ws/kernel"
)

const (
	Conversation  = 1
	Login         = 2
	Logout        = 3
	GetOnlineInfo = 4
)

type wsPipeLineFn func([]byte, *kernel.Connection) ([]byte, error)
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//推送格式
type PushData struct {
	EventType      int32       `json:"event_type"`
	PublishAccount []string    `json:"publish_account"`
	Data           interface{} `json:"data"`
}

func (response Response) Json(msg string, code int, data interface{}) []byte {
	response.Msg = msg
	response.Code = code
	response.Data = data
	res, err := json.Marshal(response)
	if err != nil {
		return []byte(`["code":404,"msg":"数据错误","data":""]`)
	}
	return res
}

//格式转换
func (pushData PushData) ConversionJson() string {
	data := pushData.Data
	switch data.(type) {
	case string:
		return data.(string)
	case float64:
		return strings.TrimRight(strconv.FormatFloat(data.(float64), 'E', -1, 64), `E+00`)
	case types.Slice,types.Map:
		result, err := json.Marshal(data)
		if err != nil {
			println(err.Error())
			return ""
		}
		return string(result)

	default:
		return ""
	}
}

//转发
func (pushData PushData) messageForwarding() {
	for _, publishAccount := range pushData.PublishAccount {
		if node, ok := kernel.GetNode(publishAccount); ok {
			go func() {
				if err := node.Ws.WriteMsg([]byte(pushData.ConversionJson())); err != nil {
					log.Println("data from ws: ", publishAccount, ":", err.Error())
				}
			}()
		}
	}
}
