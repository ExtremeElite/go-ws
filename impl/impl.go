package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"ws/conf"
)
type Connection struct {
	WsConn *websocket.Conn
	readChan chan []byte
	writeChan chan []byte
	closeChan chan byte
	one sync.Once
	isClose bool
}

func BuildConn(wsConn *websocket.Conn)(conn *Connection,err error){
	var common=conf.Config().Common
	var writeChan=common.WriteChan
	var readChan=common.ReadChan
	conn=&Connection{
		WsConn:    wsConn,
		readChan:  make(chan []byte,writeChan),
		writeChan: make(chan []byte,readChan),
		closeChan: make(chan byte,1),
	}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMsg() (data []byte,err error){
	select {
	case data=<-conn.readChan:
	case <-conn.closeChan:
		err=errors.New("读连接关闭")
	}
	return
}

func (conn *Connection) WriteMsg(data []byte) (err error)  {
	select {
	case <-conn.closeChan:
		err=errors.New("写连接关闭")
	case conn.writeChan<-data:
	}
	return
}

func (conn *Connection) Close(){
	conn.one.Do(func() {
		conn.WsConn.Close()
		conn.isClose=true
		close(conn.closeChan)
	})
}

func (conn *Connection) readLoop(){
	var (
		data []byte
		err error
	)
	for  {
		if _,data,err=conn.WsConn.ReadMessage();err!=nil {
			goto Err
		}
		select {
			case conn.readChan<-data:
			case <-conn.closeChan:
				goto Err
			
		}
	}
	Err:
		conn.Close()
}

func (conn *Connection) writeLoop()  {
	var(
		data []byte
		err error
		isClose bool
	)
	for{
		select {
		case data,isClose=<-conn.writeChan:
			if !isClose {
				goto Err
			}
		case <-conn.closeChan:
			goto Err
		}

		if err=conn.WsConn.WriteMessage(websocket.TextMessage,data);err!=nil {
			goto Err
		}
	}
	Err:
		conn.Close()
}