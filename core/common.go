package core

import "time"
const (
	HELLO=`<h1>欢迎来到Gorouting即时通讯服务</h1>`
)
func (conn *Connection) Ping() {
	go func() {
		for{
			select {
			case <-conn.closeChan:
				return
			default:
				conn.WriteMsg([]byte("Pong"))

			}
			time.Sleep(60*time.Second)
		}
	}()
}

func (conn *Connection) Pong(data []byte){
	dataString:=string(data)
	if dataString=="Ping" {
		conn.WriteMsg([]byte("Pong"))
	}
}
