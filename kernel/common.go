package kernel

import (
	"log"
	"strings"
	"time"
	"ws/common"
)

func (conn *Connection) Ping(message []byte, _ *Connection) (data []byte, err error) {
	data = message
	if strings.ToLower(string(message)) == `ping` {
		if err = conn.WriteMsg([]byte(`Pong`)); err != nil {
			common.LogDebug("写入失败:" + err.Error())
		}
		data = nil
		return
	}
	if strings.ToLower(string(message)) == `pong` {
		if err = conn.WriteMsg([]byte(`Ping`)); err != nil {
			common.LogDebug("写入失败:" + err.Error())
		}
		data = nil
		return
	}

	return
}

func (conn *Connection) Pong() {
	wsTimeOut := common.Conf.WebSocket.WsTimeOut
	if wsTimeOut <= 0 {
		return
	}

	addTime := time.Duration(wsTimeOut-1) * time.Second
	ticker := time.NewTicker(addTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if conn.IsClosed() {
				return
			}
			if err := conn.WriteMsg([]byte(`Pong`)); err != nil {
				log.Printf("[failed] auto pong: %v", err)
				return
			}
			common.LogDebug("客户端没有ping服务器,自动发送pong")
		case <-conn.closeChan:
			return
		}
	}
}
