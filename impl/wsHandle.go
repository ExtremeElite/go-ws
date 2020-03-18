package impl

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	TIME_OUT=60*time.Second
)
var upgrader=websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func WsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		dataJson string
		wsConn *websocket.Conn
		conn *Connection
		message []byte
	)
	//普通 HTTP请求
	if r.Header.Get("Connection")!="Upgrade" {
		dataJson,err=HttpAuth(r)
		if err!=nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(dataJson))
		return
	}
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		return
	}

	if conn,err= BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		goto Err
	}
	//登录判断
	if err=WsAuth(r);err!=nil {
		log.Println(err.Error())
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		goto Err
	}
	go func() {
		for{
			select {
			case data:=<-HttpChan:
				conn.WriteMsg(data)
				fmt.Println("http推送",string(data))
			}
		}
	}()
	for {
		//超时设置
		conn.WsConn.SetReadDeadline(time.Now().Add(TIME_OUT))
		if message, err=conn.ReadMsg() ;err!=nil{
			log.Println("read:", err)
			goto Err
		}
		conn.WriteMsg(message)
		log.Printf("recv: %s", message)
	}
Err:
	conn.Close()
}
