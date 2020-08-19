package core

import (
	"log"
	"net/http"
	"time"
	"ws/conf"
)

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
