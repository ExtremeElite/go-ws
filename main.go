package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/broker"
	"ws/conf"
	"ws/core"
	"ws/db"
	"ws/router"
)

func init()  {
	broker.HttpChan =make(chan broker.PushData,1)
}
func main() {
	defer func() {
		rawDB,err:=db.DB.DB()
		rawDB.Close()
		log.Println(err.Error())
	}()
	go httpPush()
	go getDataFromHttp()
	wsPush()
}
func httpPush() {
	var httpPort=conf.CommonSet.HttpPort
	var httpTimeOut=conf.CommonSet.HttpTimeOut
	httpPush:=http.NewServeMux()

	httpPush.HandleFunc("/", router.HttpRouter())
	httpPushTimeOut:=http.TimeoutHandler(httpPush,time.Duration(httpTimeOut)*time.Second,"请求超时")
	log.Printf("http服务器127.0.0.1:%d",httpPort)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPushTimeOut);err!=nil{
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort= conf.CommonSet.WsPort
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", router.WsRouter())
	log.Printf("ws服务器127.0.0.1:%d",wsPort)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal("main:",err)
	}
}
func getDataFromHttp()  {
	for{
		select {
		case data:=<-broker.HttpChan:
			core.Nodes.Range(func(name, node interface{}) bool {
				go func() {
					if err:=node.(*core.Node).Ws.WriteMsg([]byte(data.Data));err!=nil{
						log.Println("data from http: ",err.Error())
					}
				}()
				return true
			})
			log.Println(`收到的http请求推送内容:`+data.Data)
		}
	}
}

