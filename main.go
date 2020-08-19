package main

import (
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/core"
	"ws/middleware"
	"ws/middleware/auth"
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
	httpPush:=http.NewServeMux()
	httpPush.HandleFunc("/", middleware.Use(
		core.HttpHandle,middleware.Logging(),
		middleware.Method("GET"),
		auth.HttpAuthMiddle(),
		))
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPush);err!=nil{
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
					node.(*core.Node).Ws.WriteMsg(data)
				}()
				return true
			})
			log.Println(`收到的http请求推送内容:`+string(data))
		}
	}
}

