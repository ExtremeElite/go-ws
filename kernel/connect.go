package kernel

import (
	"errors"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	"ws/common"
	"ws/util"

	"github.com/gorilla/websocket"
)

var messageType = websocket.TextMessage

var panicCount int64

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
	isClose   atomic.Bool
}

func BuildConn(wsConn *websocket.Conn) (*Connection, error) {
	c := &Connection{
		WsConn:    wsConn,
		readChan:  make(chan []byte, common.Conf.WebSocket.ReadChan),
		writeChan: make(chan []byte, common.Conf.WebSocket.WriteChan),
		closeChan: make(chan struct{}),
	}
	go safeGo("readLoop", c.readLoop)
	go safeGo("writeLoop", c.writeLoop)
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
		c.isClose.Store(true)
		close(c.closeChan)
		c.WsConn.Close()
	})
}

func (c *Connection) IsClosed() bool {
	return c.isClose.Load()
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.WsConn.SetReadDeadline(t)
}

func (c *Connection) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&panicCount, 1)
			log.Printf("[PANIC] readLoop: %v", r)
			if common.Conf.Panic.LogStack {
				log.Printf("[PANIC] stack: %s", debug.Stack())
			}
		}
		c.Close()
	}()

	for {
		_, data, err := c.WsConn.ReadMessage()
		if err != nil {
			return
		}
		select {
		case c.readChan <- data:
		case <-c.closeChan:
			return
		}
	}
}

func (c *Connection) writeLoop() {
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&panicCount, 1)
			log.Printf("[PANIC] writeLoop: %v", r)
			if common.Conf.Panic.LogStack {
				log.Printf("[PANIC] stack: %s", debug.Stack())
			}
		}
		c.Close()
	}()

	for {
		select {
		case data := <-c.writeChan:
			if err := c.WsConn.WriteMessage(messageType, data); err != nil {
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func safeGo(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&panicCount, 1)
			log.Printf("[PANIC] goroutine %s: %v", name, r)
			if common.Conf.Panic.LogStack {
				log.Printf("[PANIC] stack: %s", debug.Stack())
			}
		}
	}()
	fn()
}

func GetPanicCount() int64 {
	return atomic.LoadInt64(&panicCount)
}
