package kernel

import (
	"strings"
	"time"
	"ws/common"
)

//客户端主动ping服务器自动返回
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
	var wsTimeOut = common.Setting.WsTimeOut
	if wsTimeOut > 0 {
		addTime := time.Duration(wsTimeOut-1) * time.Second
		timer := time.NewTimer(addTime)
		for range timer.C {
			if conn.IsClose {
				timer.Stop()
				goto Over
			}
			_ = conn.WriteMsg([]byte(`Pong`))
			common.LogDebug("Pong")
			timer.Reset(addTime)
		}
	}
Over:
}
