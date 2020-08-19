package core

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"ws/conf"
)
var upgrader=websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func WsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		conn *Connection
		name string
	)
	log.Println("客户端连接地址:",r.RemoteAddr)
	//普通 HTTP请求
	if r.Header.Get("Connection")!="Upgrade" {
		w.Write([]byte(HELLO))
		return
	}
	if conn,name,err=wsRequest(w,r);err!=nil {
		log.Println("wsRequest:", err.Error())
		return
	}
	AddNode(&Node{conn,name})
	for {
		//超时设置
		if err=wsRequestDone(conn);err!=nil {
			goto Err
		}
	}
	Err:
		if !conn.isClose {
			DelNode(name)
			conn.Close()
		}
}

func wsRequestDone(conn *Connection ) (err error)  {
	var message []byte
	var wsTimeOut=conf.Config().Common.WsTimeOut
	conn.WsConn.SetReadDeadline(time.Now().Add(time.Duration(wsTimeOut)*time.Second))
	if message, err=conn.ReadMsg() ;err!=nil{
		log.Println("读:", err.Error())
		return 
	}
	if err=conn.WriteMsg(message);err!=nil {
		log.Println("写:", err.Error())
		return
	}
	log.Printf("服务器收到的: %s\n", message)
	return
}
//创建连接
func wsRequest(w http.ResponseWriter,r *http.Request)(conn *Connection,name string,err error)  {
	var(
		wsConn *websocket.Conn
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		return
	}
	if conn,err= BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		conn.Close()
		return
	}
	return
}