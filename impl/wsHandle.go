package impl

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Node struct {
	ws *Connection
	name string
	id int
}
var NodeList map[string]Node
var mut sync.Mutex
const (
	TimeOut =60*time.Second
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
		name string
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		return
	}
	if conn,err= BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		return 
	}
	//登录判断
	if name,err=WsAuth(r);err!=nil {
		log.Println(err.Error())
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		return
	}
	mut.Lock()
	//关闭原来的连接
	if node,ok:=NodeList[name];ok{
		node.ws.WsConn.WriteMessage(websocket.TextMessage,[]byte(`你的连接已经在其他地方重新连接`))
		err=errors.New(`你的连接已经在其他地方重新连接`)
		node.ws.Close()
		delete(NodeList,name)
	}
	NodeList[name]=Node{
		ws:   conn,
		name: name,
		id: len(NodeList)+1,
	}
	mut.Unlock()
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
	conn.WsConn.SetReadDeadline(time.Now().Add(TimeOut))
	if message, err=conn.ReadMsg() ;err!=nil{
		log.Println("read:", err)
	}
	conn.WriteMsg(message)
	log.Printf("recv: %s", message)
	return
}
