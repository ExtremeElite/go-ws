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
		conn *Connection
	)
	//普通 HTTP请求
	if err=httpRequest(w,r);err!=nil {
		return
	}
	if conn,err=wsRequest(w,r);err!=nil {
		return
	}
	go httpOtherRequest(conn)
	for {
		//超时设置
		if err=wsRequestDone(conn);err!=nil {
			goto Err
		}
	}
Err:
	conn.Close()
}
//ws地址的普通http请求 包括数据验证
func httpRequest(w http.ResponseWriter,r *http.Request) (err error) {
	var (
		dataJson string
	)
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
func wsRequest(w http.ResponseWriter,r *http.Request)(conn *Connection,err error)  {
	var(
		wsConn *websocket.Conn
	)
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
Err:
	conn.Close()

return
}
//来自其他http服务器的数据监听
func httpOtherRequest(conn *Connection)  {
	for{
		select {
		case data:=<-HttpChan:
			conn.WriteMsg(data)
			fmt.Println("http推送",string(data))
		}
	}
}
//数据验证通过之后的数据处理部分
func wsRequestDone(conn *Connection ) (err error)  {
	var message []byte
	conn.WsConn.SetReadDeadline(time.Now().Add(TIME_OUT))
	if message, err=conn.ReadMsg() ;err!=nil{
		log.Println("read:", err)
	}
	conn.WriteMsg(message)
	log.Printf("recv: %s", message)
	return
}
