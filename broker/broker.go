/*
date:2020\11\27 0027 22:02
email:gorouting@qq.com
author:gorouting
description:
*/
package broker

import (
	"encoding/json"
	"fmt"
	"go/types"
	"strconv"
	"strings"
	"ws/common"
	"ws/kernel"
	"ws/util"
)

const (
	Conversation         = 1 //信息转发
	Login                = 2
	Logout               = 3
	GetOnlineInfo        = 4
	MessageHeaderSuccess = `[success]`
	MessageHeaderFailed  = `[failed]`
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

const NotFound = `["code":404,"msg":"数据错误","data":""]`

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

//格式转换
func (pushData PushData) ConversionJson() string {
	data := pushData.Data
	switch v := data.(type) {
	case string:
		return v
	case float64:
		return strings.TrimRight(strconv.FormatFloat(v, 'E', -1, 64), `E+0`)
	case types.Slice, types.Map:
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
		node, ok := kernel.GetNode(publishAccount)
		if ok {
			go util.Go(func() {
				if err := node.Ws.WriteMsg([]byte(pushData.ConversionJson())); err != nil {
					common.LogDebug(fmt.Sprintf("node is %s,remote_ip:%s,%s data from ws failed:%s ", node.Name, node.RemoteAddr, publishAccount, err.Error()))
				}
			})
		}

	}
}
