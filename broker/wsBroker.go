/**
 * @date:2020\11\27 0027 22:03
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package broker

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	"ws/conf"
	"ws/core"
)

func wsBroker(conn *core.Connection) (err error) {
	var message []byte
	if message, err = conn.ReadMsg(); err != nil {
		log.Println("读取失败:", err.Error())
		return
	}
	//读取消息并且发送消息
	log.Printf("服务器收到的: %s\n", message)
	if err = sendMessage(message, conn,
		ping,
		wsMessageForwarding,
	); err != nil {
		log.Printf("服务器转发消息失败: %s\n", message)
		log.Printf("错误消息引起的原因: %s\n", err.Error())
	}
	return
}
func sendMessage(message []byte, conn *core.Connection, callback ...wsPipeLineFn) error {
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

//客户端主动ping服务器
func ping(message []byte, conn *core.Connection) (data []byte, err error) {
	data = message
	if strings.ToLower(string(message)) == `ping` {
		if err = conn.WriteMsg([]byte(`Pong`)); err != nil {
			log.Println("写入失败:", err.Error())
		}
		data = nil
		return
	}
	if strings.ToLower(string(message)) == `pong` {
		if err = conn.WriteMsg([]byte(`Ping`)); err != nil {
			log.Println("写入失败:", err.Error())
		}
		data = nil
		return
	}

	return
}

//ws 消息转发 todo ws消息转发需要对连接权限进行认证
func wsMessageForwarding(message []byte, conn *core.Connection) (data []byte, err error) {
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

//go程主动pong客户端
func pong(conn *core.Connection) {
	var wsTimeOut = conf.CommonSet.WsTimeOut
	if wsTimeOut > 0 {
		addTime := time.Duration(wsTimeOut-1) * time.Second
		timer := time.NewTimer(addTime)
		for {
			select {
			case <-timer.C:
				if conn.IsClose {
					timer.Stop()
					goto Over
				}
				conn.WriteMsg([]byte(`Pong`))
				log.Println(`Pong`)

			}
			timer.Reset(addTime)
		}
	}
Over:
}
