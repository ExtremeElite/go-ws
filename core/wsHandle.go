package core

import (
	"fmt"
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
		dataJson,name string
	)
	fmt.Println(r.RemoteAddr)
	//普通 HTTP请求
	if dataJson,err=httpRequest(w,r);(err!=nil || len(dataJson)!=0){
		return
	}
	if conn,name,err=wsRequest(w,r);err!=nil||len(name)==0 {
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
//ws地址的普通http请求 包括数据验证
func httpRequest(w http.ResponseWriter,r *http.Request) (dataJson string,err error) {
	if r.Header.Get("Connection")!="Upgrade" {
		dataJson,err=HttpAuth(r)
		if err!=nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(dataJson))
	}
	return
}
//ws地址的ws请求 包括数据验证
func wsRequest(w http.ResponseWriter,r *http.Request)(conn *Connection,name string,err error)  {
	var(
		wsConn *websocket.Conn
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		return
	}
	if conn,err= BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		conn.Close()
		return
	}
	//登录判断
	if name,err=WsAuth(r);err!=nil {
		log.Println(err.Error())
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		conn.Close()
	}
	return
}
//数据验证通过之后的数据处理部分
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
