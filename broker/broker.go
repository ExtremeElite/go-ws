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
	Conversation  = 1 //信息转发
	Login         = 2
	Logout        = 3
	GetOnlineInfo = 4
)

type wsPipeLineFn func([]byte, *kernel.Connection) ([]byte, error)

//推送格式
type PushData struct {
	EventType      int32       `json:"event_type"`
	PublishAccount []string    `json:"publish_account"`
	Data           interface{} `json:"data"`
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
		if node, ok := kernel.GetNode(publishAccount); ok {
			go util.Go(func() {
				if err := node.Ws.WriteMsg([]byte(pushData.ConversionJson())); err != nil {
					common.LogDebug(fmt.Sprintf("node is %s,remote_ip:%s,%s data from ws failed:%s ", node.Name, node.RemoteAddr, publishAccount, err.Error()))
				}
			})
		} else {
			common.LogDebug("messageForeard node not found")
		}
	}
}
