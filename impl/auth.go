package impl

import (
	"time"
)

var(
	conn *Connection
	pongNum uint8
	data []byte
	err error
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
			time.Sleep(30*time.Second)
		}
	}()
}

func (conn *Connection) Ping()  {
	go func() {
		for{
			if data,err=conn.ReadMsg();err!=nil{
				goto Err
			}
			dataString:=string(data)
			if dataString=="ping" {

			}
		}
		Err:
			conn.Close()
	}()
}
