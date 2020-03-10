package impl

import (
	"time"
)

var(
	conn *Connection
	pongNum uint8
	data []byte
	err error
	dataString string
	currentTime int64
)
func (conn *Connection) IsAuth()  {
	
}
func (conn *Connection) Pong() {
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

func (conn *Connection) Ping(data []byte)  {
	dataString=string(data)
	if currentTime==0 || time.Now().Unix()-currentTime<60*2 {
		currentTime=time.Now().Unix()
		conn.WriteMsg([]byte("Pong"))
		return
	}
	conn.Close()
}
