package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/conf"
	"ws/core"
	"ws/pipeLine"
)

func main() {
	run()
}

func run(){
	core.HttpChan=make(chan []byte,1)
	go httpPush()
	go getDataFromHttp()
	wsPush()
}
func httpPush() {
	var httpPort=conf.CommonSet.HttpPort
	var httpTimeOut=conf.CommonSet.HttpTimeOut
	httpPush:=http.NewServeMux()

	httpPush.HandleFunc("/", pipeLine.Use(
		core.HttpHandle,
		pipeLine.Logging(),
		pipeLine.Method("GET"),
		pipeLine.HttpAuthMiddle(),
		))
	httpPushTimeOut:=http.TimeoutHandler(httpPush,time.Duration(httpTimeOut)*time.Second,"请求超时")
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPushTimeOut);err!=nil{
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort= conf.CommonSet.WsPort
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", core.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal("main:",err)
	}
}
func getDataFromHttp()  {
	for{
		select {
		case data:=<-core.HttpChan:
			core.Nodes.Range(func(name, node interface{}) bool {
				go func() {
					if err:=node.(*core.Node).Ws.WriteMsg(data);err!=nil{
						log.Println("data from http: ",err.Error())
					}
				}()
				return true
			})
			log.Println(`收到的http请求推送内容:`+string(data))
		}
	}
}

