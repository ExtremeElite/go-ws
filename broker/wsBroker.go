/*
date:2020\11\27 0027 22:03
email:gorouting@qq.com
author:gorouting
description:
*/
package broker

import (
	"encoding/json"
	"errors"
	"ws/common"
	"ws/kernel"
	"ws/util"
)

func wsBroker(conn *kernel.Connection) (err error) {
	var message []byte
	if message, err = conn.ReadMsg(); err != nil {
		common.LogInfoFailed("读取失败:" + err.Error())
		return
	}
	//读取消息并且发送消息
	common.LogInfoSuccess("服务器收到的消息:" + string(message))
	if err = sendMessage(message, conn,
		conn.Ping,
		wsMessageForwarding,
	); err != nil {
		//将错误发送给客户端
		var returnClientMsg = util.Response{}
		var byteReturnClientMsg []byte = returnClientMsg.Json(err.Error(), 404, string(message))
		conn.WriteMsg(byteReturnClientMsg)
		//服务器埋点
		common.LogInfoFailed("服务器转发消息:" + string(byteReturnClientMsg))
	}
	return
}
func sendMessage(message []byte, conn *kernel.Connection, callback ...wsPipeLineFn) error {
	_message := message
	var err error
	for _, fn := range callback {
		if _message != nil {
			if _message, err = fn(_message, conn); err != nil {
				return err
			}
		}
	}
	return nil
}

// ws 消息转发 todo ws消息转发需要对连接权限进行认证
func wsMessageForwarding(message []byte, conn *kernel.Connection) (data []byte, err error) {
	var pushData PushData
	data = message
	if err = json.Unmarshal(message, &pushData); err != nil {
		err = errors.New(`wsMessageForwarding data error:` + err.Error())
		return
	}
	switch pushData.EventType {
	case Conversation:
		pushData.messageForwarding()
	}
	return
}
