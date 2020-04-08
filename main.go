package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/core"
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
	httpPush.HandleFunc("/", core.HttpHandle)
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
			core.NodeList.Range(func(name, node interface{}) bool {
				node.(*core.Node).Ws.WriteMsg(data)
				return true
			})
			fmt.Println(string(data))
		}
	}
}

