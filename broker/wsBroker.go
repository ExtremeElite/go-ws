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
)

func wsBroker(conn *kernel.Connection) (err error) {
	var message []byte
	if message, err = conn.ReadMsg(); err != nil {
		common.LogDebug("读取失败:" + err.Error())
		return
	}
	//读取消息并且发送消息
	common.LogInfo("服务器收到的:\n" + string(message))
	if err = sendMessage(message, conn,
		conn.Ping,
		wsMessageForwarding,
	); err != nil {
		common.LogDebug("服务器转发消息失败:\n" + string(message) + "错误消息引起的原因:\n" + err.Error())
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

//ws 消息转发 todo ws消息转发需要对连接权限进行认证
func wsMessageForwarding(message []byte, conn *kernel.Connection) (data []byte, err error) {
	var pushData PushData
	data = message
	if err = json.Unmarshal(message, &pushData); err != nil {
		err = errors.New(`wsMessageForwarding data error`)
		return
	}
	switch pushData.EventType {
	case Conversation:
		pushData.messageForwarding()
	}
	return
}
