package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	WsConn *websocket.Conn
	ReadChan chan []byte
	WriteChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClose bool
}

func BuildConn(wsConn *websocket.Conn)(conn *Connection,err error){
	conn=&Connection{
		WsConn:    wsConn,
		ReadChan:  make(chan []byte,1000),
		WriteChan: make(chan []byte,1000),
		closeChan: make(chan byte,1),
	}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

func (ws *Connection) ReadMsg() (data []byte,err error){
	select {
	case data=<-ws.ReadChan:
	case <-ws.closeChan:
		err=errors.New("connected is closed")
	}
	return
}

func (ws *Connection) WriteMsg(data []byte) (err error)  {
	select {
	case ws.WriteChan<-data:
	case <-ws.closeChan:
		err=errors.New("connected is closed")
	}

	return
}

func (ws *Connection) Close(){
	ws.WsConn.Close()

	ws.mutex.Lock()
	if !ws.isClose {
		ws.isClose=true
		close(ws.closeChan)
	}
	ws.mutex.Unlock()
}

func (ws *Connection) readLoop(){
	var (
		data []byte
		err error
	)
	for  {
		if _,data,err=ws.WsConn.ReadMessage();err!=nil {
			goto Err
		}
		select {
			case ws.ReadChan<-data:
			case <-ws.closeChan:
				goto Err
			
		}
	}
	Err:
		ws.Close()
}

func (ws *Connection) writeLoop()  {
	var(
		data []byte
		err error
	)
	for{
		select {
		case data=<-ws.WriteChan:
		case <-ws.closeChan:
			goto Err
		}

		if err=ws.WsConn.WriteMessage(websocket.TextMessage,data);err!=nil {
			goto Err
		}
	}
	Err:
		ws.Close()
}