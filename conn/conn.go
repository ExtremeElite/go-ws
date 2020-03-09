package conn

import "github.com/gorilla/websocket"

type Connection struct {
	WsConn *websocket.Conn
	ReadChan chan []byte
	WriteChan chan []byte
}

func BuildConn(wsConn *websocket.Conn)(conn *Connection,err error){
	conn=&Connection{
		WsConn:    wsConn,
		ReadChan:  make(chan []byte,1000),
		WriteChan: make(chan []byte,1000),
	}
	return
}

func (ws *Connection) ReadMsg(conn websocket.Conn) (err error){
	var (
		msg []byte
	)
	_,msg,err=conn.ReadMessage()
	ws.ReadChan<-msg
	return
}

func (ws * Connection) WriteMsg(conn websocket.Conn) (err error)  {
	var (
		msg []byte
	)
	msg=<-ws.WriteChan
	err=conn.WriteMessage(websocket.TextMessage,msg)
	return
}