package core

import "time"

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
	dataString=string(data)
	if dataString=="Ping" {
		conn.WriteMsg([]byte("Pong"))
	}
}
