package kernel

import (
	"errors"
	"log"
	"sync"
	"time"
	"ws/common"
	"ws/util"

	"github.com/gorilla/websocket"
)

var messageType = websocket.TextMessage

func init() {
	if common.Conf.MessageType != 1 {
		messageType = websocket.BinaryMessage
	}
}

type Connection struct {
	WsConn    *websocket.Conn
	readChan  chan []byte
	writeChan chan []byte
	closeChan chan struct{}
	once      sync.Once
	IsClose   bool
}

func BuildConn(wsConn *websocket.Conn) (*Connection, error) {
	c := &Connection{
		WsConn:    wsConn,
		readChan:  make(chan []byte, common.Conf.WebSocket.ReadChan),
		writeChan: make(chan []byte, common.Conf.WebSocket.WriteChan),
		closeChan: make(chan struct{}),
	}
	go safeGo(c.readLoop)
	go safeGo(c.writeLoop)
	return c, nil
}

func (c *Connection) ReadMsg() ([]byte, error) {
	select {
	case data := <-c.readChan:
		return data, nil
	case <-c.closeChan:
		return nil, errors.New(util.ReadConnectClosed)
	}
}

func (c *Connection) WriteMsg(data []byte) error {
	select {
	case <-c.closeChan:
		return errors.New(util.WriteConnectClosed)
	case c.writeChan <- data:
		return nil
	}
}

func (c *Connection) Close() {
	c.once.Do(func() {
		c.IsClose = true
		close(c.closeChan)
		c.WsConn.Close()
	})
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.WsConn.SetReadDeadline(t)
}

func (c *Connection) readLoop() {
	for {
		_, data, err := c.WsConn.ReadMessage()
		if err != nil {
			c.Close()
			return
		}
		select {
		case c.readChan <- data:
		case <-c.closeChan:
			c.Close()
			return
		}
	}
}

func (c *Connection) writeLoop() {
	for {
		select {
		case data := <-c.writeChan:
			if err := c.WsConn.WriteMessage(messageType, data); err != nil {
				c.Close()
				return
			}
		case <-c.closeChan:
			c.Close()
			return
		}
	}
}

func safeGo(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("goroutine panic:", r)
		}
	}()
	fn()
}
